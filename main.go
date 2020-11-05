package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	dbw "github.com/DanielRHolland/qrep/db"
	. "github.com/DanielRHolland/qrep/models"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/mux"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type appContext struct {
	addr          string
	base          string
	db            dbw.DbAccessor
	tokenUsers    map[string]string
	userPasswords map[string]string
}

func newAppContext(addr string, base string, usersFile string, dbConnection dbw.DbAccessor) *appContext {
	if base == "" {
		base = "http://" + addr
	}
	return &appContext{
		addr:          addr,
		base:          base,
		db:            dbConnection,
		tokenUsers:    make(map[string]string),
		userPasswords: *loadUsers(usersFile),
	}
}

//Placeholder for more sophisticated auth system
func loadUsers(filename string) *map[string]string {
	if filename == "" {
		return &map[string]string{"admin": "admin"}
	} else {
		readFile, err := os.Open(filename)
		if err != nil {
			log.Fatal("Could not read users file: ", filename)
		}
		result := make(map[string]string)

		fileScanner := bufio.NewScanner(readFile)
		fileScanner.Split(bufio.ScanLines)
		for fileScanner.Scan() {
			temp := strings.Split(fileScanner.Text(), ",")
			result[temp[0]] = temp[1]
		}
		readFile.Close()
		return &result
	}
}

//POST create new qrcode
func (ctx *appContext) createQr(w http.ResponseWriter, r *http.Request) {
	var item TrackedItemType
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		checkError(err)
		item.Name = r.Form.Get("name")
	} else if r.Header.Get("Content-Type") == "application/json" {
		log.Println(r.Header)
		reqBody, err := ioutil.ReadAll(r.Body)
		checkError(err)
		json.Unmarshal(reqBody, &item)
	}
	if item.Issues == nil {
		item.Issues = []IssueType{}
	}
	item.Id, _ = ctx.db.InsertItem(item)
	http.Redirect(w, r, "", 303)
}

//GET qrCreation page
func serveCreationPage(w http.ResponseWriter, r *http.Request) {
	renderCreationPage(w)
}

//GET qr code page for id
func (ctx *appContext) serveQr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if item, err := ctx.db.GetItem(id); err == nil { // qr
		renderQr(w, item)
	} else {
		io.WriteString(w, "NOPE")
	}
}

//GET qr code png for id

func (ctx *appContext) serveQrPng(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx.encodeQrPng(w, vars["id"])
}

func (ctx *appContext) encodeQrPng(w io.Writer, id string) {
	url := ctx.base + "/" + id
	qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	png.Encode(w, qrCode)
}

//GET reporting page for item
func serveReportingPage(w http.ResponseWriter, _ *http.Request) {
	renderReportingPage(w)
}

//POST new report
func (ctx *appContext) newReportPosted(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //should err check here
	vars := mux.Vars(r)
	id, exists := vars["id"]
	var issue IssueType
	issue.Description = r.Form.Get("issue")
	if exists {
		ctx.db.AddIssueToItem(issue, id)
	}
	thanksForReport(w)
}

func (ctx *appContext) serveItemReportLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, varPresent := vars["id"] //returns id
	if varPresent {
		if item, err := ctx.db.GetItem(id); err == nil {
			renderItemReportLog(w, item)
		} else {
			io.WriteString(w, "NOPE")
		}
	}
}

func (ctx *appContext) serveReportLog(w http.ResponseWriter, _ *http.Request) {
	trackedItems, err := ctx.db.GetTrackedItems(100) //get tracked items
	if err == nil {
		renderReportLog(w, trackedItems)
	} else {
		io.WriteString(w, "NOPE")
	}
}

func (ctx *appContext) serveDashboard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, namePresent := query["name"]
	var trackedItems []TrackedItemType
	var err error
	if !namePresent || len(name) <= 0 || name[0] == "" {
		trackedItems, err = ctx.db.GetTrackedItems(100)
	} else {
		trackedItems, err = ctx.db.SearchTrackedItems(100, name[0])
	}
	if err == nil {
		renderDashboard(w, trackedItems)
	} else {
		io.WriteString(w, "NOPE")
	}
}

func (ctx *appContext) updateIssue(w http.ResponseWriter, r *http.Request) {
	var issue IssueType
	reqBody, err := ioutil.ReadAll(r.Body)
	checkError(err)
	json.Unmarshal(reqBody, &issue)
	ctx.db.UpdateDbIssue(issue)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func (ctx *appContext) serveItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	itemids, itemidsPresent := query["item"]
	if itemidsPresent && len(itemids) > 0 {
		items, _ := ctx.db.GetItemsFromIds(itemids)

		action, actionPresent := query["action"]
		if actionPresent && len(action) > 0 {
			switch action[0] {
			case "getqrs":
				index := rand.Int()
				pdfcache[index] = parGenQrPdf(items, ctx.base)
				indexstring := strconv.Itoa(index)
				renderItemsQrsPage(w, indexstring)
			case "getqrszip":
				ctx.generateQrsZip(w, items)
			}
		}
	}
}

func parGenQrPdf(items []TrackedItemType, baseurl string) chan bytes.Buffer {
	pdfChan := make(chan bytes.Buffer)
	go func() { pdfChan <- generateQrsPdf(items, baseurl) }()
	return pdfChan
}

var pdfcache = make(map[int]chan bytes.Buffer)

func serveCachedPdf(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	index, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	if pdfChan, exists := pdfcache[index]; exists {
		pdfbytes := <-pdfChan
		pdfbytes.WriteTo(w)
		close(pdfChan)
		delete(pdfcache, index)
	}
}

func (ctx *appContext) generateQrsZip(w http.ResponseWriter, items []TrackedItemType) {
	zipWriter := zip.NewWriter(w)
	for _, v := range items {
		pngWriter, _ := zipWriter.Create(v.Name + "-" + v.Id + ".png")
		ctx.encodeQrPng(pngWriter, v.Id)
	}
	zipWriter.Close() //errors to handle
}

func (ctx *appContext) removeItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	itemids, itemidsPresent := query["item"]
	if itemidsPresent && len(itemids) > 0 {
		ctx.db.RemoveItemsFromDb(itemids)
	}
	http.Redirect(w, r, "", 303)
}

func (ctx *appContext) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderLoginPage(w)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		checkError(err)
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		if ctx.checkUserPasswordValid(username, password) {
			expire := time.Now().Add(30 * time.Minute)
			token := "key"
			cookie := http.Cookie{
				Name:    "sessionToken",
				Value:   token,
				Expires: expire,
			}
			ctx.tokenUsers[token] = username
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", 303)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	} else {
		http.Error(w, "Method not found", http.StatusNotFound)
	}
}

func (ctx *appContext) checkUserPasswordValid(username string, password string) bool {
	//TODO implement properly
	pw, exists := ctx.userPasswords[username]
	return exists && pw == password
}

// Route declaration
func (ctx *appContext) router() *mux.Router {
	r := mux.NewRouter()
	staticServer := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticServer))
	s := r.PathPrefix("").Subrouter() //.Headers("X-Session-Token").Subrouter()
	ctx.secureZoneSubRouter(s)
	r.HandleFunc("/login", ctx.login)
	r.HandleFunc("/{id}", ctx.newReportPosted).Methods("POST")
	r.HandleFunc("/{id}", serveReportingPage)
	r.Use(loggingMiddleware)
	return r
}

func (ctx *appContext) secureZoneSubRouter(r *mux.Router) { //TODO better name
	r.Use(ctx.Middleware)
	r.HandleFunc("/", ctx.serveDashboard)
	r.HandleFunc("/issue/{id}", ctx.updateIssue).Methods("PUT")
	r.HandleFunc("/new", ctx.createQr).Methods("POST")
	r.HandleFunc("/qr", ctx.createQr).Methods("POST")
	r.HandleFunc("/remove", ctx.removeItems).Methods("GET") //GET shouldn't modify server state
	r.HandleFunc("/items", ctx.serveItems).Methods("GET")
	r.HandleFunc("/dl/{id}.png", ctx.serveQrPng)
	r.HandleFunc("/dl/qrcodes.zip", ctx.serveItems)
	r.HandleFunc("/dl/{id}/qrcodes.pdf", serveCachedPdf)
	r.HandleFunc("/qr", serveCreationPage)
	r.HandleFunc("/qr/{id}", ctx.serveQr)
	r.HandleFunc("/qrpng/{id}", ctx.serveQrPng)
	r.HandleFunc("/reports/{id}", ctx.serveItemReportLog)
	r.HandleFunc("/reports", ctx.serveReportLog)
}

// Middleware function, which will be called for each request
func (ctx *appContext) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if token == "" {
			tokens := r.URL.Query()["api"]
			if len(tokens) > 0 {
				token = tokens[0]
			} else {
				cookie, err := r.Cookie("sessionToken")
				if err == nil {
					token = cookie.Value
				}
			}
		}
		if user, found := ctx.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			//	http.Error(w, "Forbidden", http.StatusForbidden)
			//ctx.login(w, r)
			http.Redirect(w, r, "login", 303)
		}
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Initiate web server
func main() {
	addr := flag.String("addr", "127.0.0.1:9100", "Server Address:port")
	base := flag.String("url", "", "URL")
	usersFile := flag.String("users", "", "File containing users and passwords")
	flag.Parse()
	dbConnection := dbw.NewMongoDbConnection()
	defer dbConnection.Disconnect()
	ctx := newAppContext(*addr, *base, *usersFile, dbConnection)
	log.Println("Serving on ", ctx.base, "(", ctx.addr, ")")
	router := ctx.router()
	srv := &http.Server{
		Handler:      router,
		Addr:         ctx.addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
