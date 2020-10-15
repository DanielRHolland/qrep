package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
)

const header = `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title>QRep</title>
            <!--Import Google Icon Font-->
            <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
            <!-- Compiled and minified CSS -->
            <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/css/materialize.min.css">
            <!--Let browser know website is optimized for mobile-->
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        </head>
        <body>
        <div class="container">
        `

const footer = `
            </div>
    </body>

    <!-- Compiled and minified JavaScript -->
    <script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0/js/materialize.min.js"></script>
</html>
`

const thanksForReport = header + "<h3>Thanks for your report!</h3>" + footer

func imageCard(title string, image string, content string) string {
	return `
            <div class="card" style="width:256px;">
              <div class="card-image waves-effect waves-block waves-light">
                <img class="activator" src="` + image + `">
              </div>
              <div class="card-content">
              <span class="card-title">` + title + `</span>
                ` + content + `
              </div>
            </div>
        `
}

//<div class="row">
//  <div class="col s12 m6 offset-m6">
//              <a href="#">` + link + `</a>
func renderQr(w http.ResponseWriter, item trackedItem) {
	var tpl = header +
		imageCard("{{.Item.Name}}", "/qrpng/{{.Item.Id}}", "<a href='/{{.Item.Id}}'>{{.Item.Id}}</a>") +
		footer
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
//	id := strconv.Itoa(number) //add error checking
	model := struct {
		Id   string
		Item trackedItem
	}{
		Item: item,
	}

	err = t.Execute(w, model)
	checkError(err)

}

// <textarea id="issue" class="materialize-textarea"></textarea>
//<input id="issue" type="text" name="issue">
func renderReportingPage(w http.ResponseWriter) {
	const page = header + `
          <h2>What's the problem?</h2>
          <form action="" method="post">
            <div class="input-field">
              <textarea id="issue" name="issue" class="materialize-textarea"></textarea>
              <label for="issue">Issue:</label>
            </div>

            <button class="btn waves-effect waves-light" type="submit" name="action">Submit
                <i class="material-icons right">send</i>
            </button>
          </form> 
       ` + footer
	io.WriteString(w, page)
}

func renderCreationPage(w http.ResponseWriter) {
	const page = header + `
        <h2>Add a new item to track</h2>
          <form action="" method="post">
            <div class="input-field">
              <input id="name" name="name" type="text"></input>
              <label for="name">Name:</label>
            </div>

            <button class="btn waves-effect waves-light" type="submit" name="action">Submit
                <i class="material-icons right">send</i>
            </button>
          </form> 
        ` + footer
	io.WriteString(w, page)
}

func renderItemReportLog(w http.ResponseWriter, item trackedItem) {
	const tpl = header + `
        <h2>Reports for {{.Name}}</h2>
        <ul>
          {{range .Issues}}
              <li> {{.}} </li>
          {{end}}
        </ul>
        ` + footer
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
	err = t.Execute(w, item)
	checkError(err)

}

func renderReportLog(w http.ResponseWriter, trackedItems []trackedItem, itemlessIssues []string) {
	const tpl = header + `
        {{range .TrackedItems}}
        <h4><a href="/reports/{{.Id}}">Reports for {{.Name}}</a></h4>
        <ul>
          {{range .Issues}}
              <li> {{.}} </li>
          {{end}}
        </ul>
        {{end}}
        {{if .ItemlessIssues}}
        <h4>Itemless issues:</h4>
        <ul>
          {{range .ItemlessIssues}}
              <li> {{.}} </li>
          {{end}}
        </ul>
        {{end}}

        ` + footer
	t, err := template.New("webpage").Parse(tpl)
	checkError(err)
	model := struct {
		TrackedItems   []trackedItem
		ItemlessIssues []string
	}{
		TrackedItems:   trackedItems,
		ItemlessIssues: itemlessIssues,
	}
	err = t.Execute(w, model)
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
