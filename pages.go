package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func renderQr(w http.ResponseWriter, item trackedItem, number int) {
	const tpl = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>{{.Item.Name}}</title>
        </head>
        <body>
            <a href="{{.Url}}">{{.Url}}</a>
            <img src="{{.QrUrl}}">
        </body>
    </html>`
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
	id := strconv.Itoa(number)
	model := struct {
		Url   template.URL
		QrUrl template.URL
		Item  trackedItem
	}{
		Url:   template.URL(base + "/" + id),
		QrUrl: template.URL(base + "/qrpng/" + id),
		Item:  item,
	}

	err = t.Execute(w, model)
	checkError(err)

}

func renderReportingPage(w http.ResponseWriter, id string) {
	const tpl = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>{{.Id}}</title>
        </head>
        <body>
          <form action="" method="post">
            <label for="issue">Issue:</label>
            <input type="text" id="issue" name="issue"><br><br>
            <input type="submit" value="Submit">
          </form> 
        </body>
    </html>`
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
	model := struct {
		Id string
	}{
		Id: id,
	}

	err = t.Execute(w, model)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
