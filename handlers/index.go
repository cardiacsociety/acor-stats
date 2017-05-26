package handlers

import (
	"html/template"
	"net/http"
	"os"
	"fmt"
)

var tpl *template.Template

func init() {

	tpl = template.Must(template.ParseGlob("./templates/*"))
}

func Index(w http.ResponseWriter, r *http.Request) {

	// Need to pass a few things to the template
	var data struct {
		APIBaseURL string
	}
	data.APIBaseURL = os.Getenv("API_BASE_URL")
	fmt.Println(data.APIBaseURL)

	tpl.ExecuteTemplate(w, "index", data)
}

func Test(w http.ResponseWriter, r *http.Request) {

	// pass the raw query to the template so it can be appended to the api url
	q := r.URL.RawQuery
	tpl.ExecuteTemplate(w, "test", q)
}

func ReportsAllHandler(w http.ResponseWriter, r *http.Request) {

}
func ReportsAllStateHandler(w http.ResponseWriter, r *http.Request) {

}
func ReportsDevicesHandler(w http.ResponseWriter, r *http.Request) {

}
func ReportsDevicesStateHandler(w http.ResponseWriter, r *http.Request) {

}
func ReportsProceduresHandler(w http.ResponseWriter, r *http.Request) {

}
func ReportsProceduresStateHandler(w http.ResponseWriter, r *http.Request) {

}
