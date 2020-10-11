package main

import (
	"encoding/json"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gorilla/mux"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const addr = "127.0.0.1:9100"
const base = "http://" + addr

var trackedItems []trackedItem
var itemlessIssues []string

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

type trackedItem struct {
	Name   string   `json:"name"`
	Issues []string `json:"issues"`
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
	key := len(trackedItems)
	trackedItems = append(trackedItems, item)
	renderQr(w, item, key)
}

//GET qrCreation page
func serveCreationPage(w http.ResponseWriter, r *http.Request) {
	renderCreationPage(w)
}

//GET qr code page for id
func serveQr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	i, _ := strconv.Atoi(key) //need err checker
	if len(trackedItems) > i && i >= 0 {
		//            io.WriteString(w, trackedItems[i].Name)
		renderQr(w, trackedItems[i], i)
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
	id := vars["id"]
	key, err := strconv.Atoi(id)
	issue := r.Form.Get("issue")
	if err != nil || len(trackedItems) <= key || key < 0 {
		itemlessIssues = append(itemlessIssues, issue)
	} else {
		issuesLog := &trackedItems[key].Issues
		*issuesLog = append(*issuesLog, issue)
	}
	io.WriteString(w, thanksForReport)
}

func serveItemReportLog(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	key, err := strconv.Atoi(id)
	checkError(err)
	if len(trackedItems) > key && key >= 0 {
		renderItemReportLog(w, trackedItems[key])
	} else {
		io.WriteString(w, "NOPE")
	}
}

func serveReportLog(w http.ResponseWriter, _ *http.Request) {
	renderReportLog(w, trackedItems, itemlessIssues)
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/qr", createQr).Methods("POST")
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
        router := router()
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
