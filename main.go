package main

import (
	"encoding/json"
	"flag"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/mux"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var addr string
var base string

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

type trackedItem struct {
	Name   string      `json:"name" bson:"name"`
	Issues []issueType `json:"issues" bson:"issues"`
	Id     string      `json:"id" bson:"_id"`
}

type issueType struct {
	Description string `json:"description" bson:"description"`
	Resolved    bool   `json:"resolved" bson:"resolved"`
	Id          string `json:"id" bson:"_id"`
}

//POST create new qrcode
func createQr(w http.ResponseWriter, r *http.Request) {
	var item trackedItem
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
		item.Issues = []issueType{}
	}
	item.Id = insertItem(item)
	http.Redirect(w, r, "", 303)
}

//GET qrCreation page
func serveCreationPage(w http.ResponseWriter, r *http.Request) {
	renderCreationPage(w)
}

//GET qr code page for id
func serveQr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if item, err := getItem(id); err == nil { // qr
		//            io.WriteString(w, trackedItems[i].Name)
		renderQr(w, item)
	} else {
		io.WriteString(w, "NOPE")
	}
}

//GET qr code png for id

func serveQrPng(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url := base + "/" + vars["id"]
	qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	png.Encode(w, qrCode)
}

//GET reporting page for item
func serveReportingPage(w http.ResponseWriter, _ *http.Request) {
	renderReportingPage(w)
}

//POST new report
func newReportPosted(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //should err check here
	vars := mux.Vars(r)
	id, exists := vars["id"]
	var issue issueType
	issue.Description = r.Form.Get("issue")

	if exists {
		//                itemExists := false
		//
		//        } else {

		//        }

		//	if exists && itemExists  { // If
		//Add new issue to the item with the id
		addIssueToItem(issue, id)
		//	} else {
		//Insert new issue into itemlessIssues
	}
	thanksForReport(w)
}

func serveItemReportLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, varPresent := vars["id"] //returns id
	if varPresent {
		if item, err := getItem(id); err == nil {
			renderItemReportLog(w, item)
		} else {
			io.WriteString(w, "NOPE")
		}
	}
}

func serveReportLog(w http.ResponseWriter, _ *http.Request) {
	trackedItems, err := getTrackedItems(100) //get tracked items
	if err == nil {
		renderReportLog(w, trackedItems)
	} else {
		io.WriteString(w, "NOPE")
	}
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name, namePresent := query["name"]
	var trackedItems []trackedItem
	var err error
	if !namePresent || len(name) <= 0 || name[0] == "" {
		trackedItems, err = getTrackedItems(100)
	} else {
		trackedItems, err = searchTrackedItems(100, name[0])
	}
	if err == nil {
		renderDashboard(w, trackedItems)
	} else {
		io.WriteString(w, "NOPE")
	}
}

func updateIssue(w http.ResponseWriter, r *http.Request) {
	var issue issueType
	reqBody, err := ioutil.ReadAll(r.Body)
	checkError(err)
	json.Unmarshal(reqBody, &issue)
	updateDbIssue(issue)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issue)
}

func serveItems(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()
        action, actionPresent := query["action"]
        if (actionPresent && len(action) > 0 && action[0] == "getqrs") {
            itemids, itemidsPresent := query["item"]
            if (itemidsPresent && len(itemids) > 0) {
                items := getItemsFromIds(itemids)
                renderItemsQrsPage(w, items)
            }
        }
        io.WriteString(w, "foo")
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()
	staticServer := http.FileServer(http.Dir("static/"))
	r.HandleFunc("/", serveDashboard)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticServer))
	r.HandleFunc("/issue/{id}", updateIssue).Methods("PUT")
	r.HandleFunc("/new", createQr).Methods("POST")
	r.HandleFunc("/qr", createQr).Methods("POST")
	r.HandleFunc("/items", serveItems)
	r.HandleFunc("/qr", serveCreationPage)
	r.HandleFunc("/qr/{id}", serveQr)
	r.HandleFunc("/qrpng/{id}", serveQrPng)
	r.HandleFunc("/reports/{id}", serveItemReportLog)
	r.HandleFunc("/reports", serveReportLog)
	r.HandleFunc("/{id}", newReportPosted).Methods("POST")
	r.HandleFunc("/{id}", serveReportingPage)
	return r
}

// Initiate web server
func main() {
	addrFlag := flag.String("addr", "127.0.0.1:9100", "Server Address:port")
	baseFlag := flag.String("url", "def", "URL")
	flag.Parse()
	addr = *addrFlag
	base = *baseFlag
	if base == "def" {
		base = "http://" + addr
	}
	log.Println("Serving on ", base, "(", addr, ")")
	router := router()
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
