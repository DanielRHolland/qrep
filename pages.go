package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template = loadTemplates()

func getTemplatePaths() []string {
	dirpath := "./templates"
	dir, err := os.Open(dirpath)
	if err != nil {
		log.Fatal(err)
	}
	filenames, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	filepaths := []string{}
	for i := range filenames {
		filename := filenames[i]
		l := len(filename)
		if filename[l-5:] == ".html" || filename[l-3:] == ".js" {
			filepaths = append(filepaths, dirpath+"/"+filename)
		}
	}
	return filepaths
}

func loadTemplates() *template.Template {
	paths := getTemplatePaths()
	t, err := template.ParseFiles(paths...)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func thanksForReport(w http.ResponseWriter) {
	if err := templates.ExecuteTemplate(w, "report_thanks.html", struct{}{}); err != nil {
		log.Fatal(err)
	}
}

func renderDashboard(w http.ResponseWriter, items []trackedItem) {
	model := struct {
		TrackedItems []trackedItem
	}{
		TrackedItems: items,
	}
	err := templates.ExecuteTemplate(w, "dashboard_page.html", model)
	if err != nil {
		log.Fatal(err)
	}
}

func renderQr(w http.ResponseWriter, item trackedItem) {
	err := templates.ExecuteTemplate(w, "qr_view.html", item)
	checkError(err)

}

func renderReportingPage(w http.ResponseWriter) {
	err := templates.ExecuteTemplate(w, "reporting_page.html", struct{}{})
	checkError(err)
}

func renderCreationPage(w http.ResponseWriter) {
	err := templates.ExecuteTemplate(w, "creation_page.html", struct{}{})
	checkError(err)
}

func renderItemReportLog(w http.ResponseWriter, item trackedItem) {
	err := templates.ExecuteTemplate(w, "item_report_log_page.html", item)
	checkError(err)

}

func renderReportLog(w http.ResponseWriter, trackedItems []trackedItem) {
	err := templates.ExecuteTemplate(w, "report_log_page.html", trackedItems)
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
