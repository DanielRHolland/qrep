package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"io"
        "io/ioutil"
        "encoding/json"
        "strconv"
        "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
        "image/png"
)

const base = "127.0.0.1:9100"

var trackedItems []trackedItem

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

type trackedItem struct {
        Name string
}

//POST create new qrcode
func createQr(w http.ResponseWriter, r *http.Request) {
        reqBody, _ := ioutil.ReadAll(r.Body)
        var item trackedItem
        json.Unmarshal(reqBody, &item)
        key := len(trackedItems)
        trackedItems = append(trackedItems, item)
        itemUri := base + "/qr/" + string(key)
        json.NewEncoder(w).Encode(itemUri)
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
        url := "http://"+base+"/"+vars["id"]
        qrCode, _ := qr.Encode(url, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 256, 256)
	png.Encode(w, qrCode)
}

//GET reporting page for item
func serveReportingPage(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        key := vars["id"]
        renderReportingPage(w, key)

}
//POST new report
func newReportPosted(w http.ResponseWriter, r *http.Request) {
        r.ParseForm() //should err check here
        log.Println(r.Form)
        io.WriteString(w, "Issue: "+r.Form.Get("issue"))
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
        r.HandleFunc("/qr", createQr).Methods("POST")
        r.HandleFunc("/qr/{id}", serveQr)
        r.HandleFunc("/qrpng/{id}", serveQrPng)
	r.HandleFunc("/{id}", newReportPosted).Methods("POST")
	r.HandleFunc("/{id}", serveReportingPage)
	return r
}

// Initiate web server
func main() {
	router := router()
	srv := &http.Server{
		Handler: router,
		Addr:    base,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
