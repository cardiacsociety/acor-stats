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

	var data struct {
		APIBaseURL string
	}
	data.APIBaseURL = os.Getenv("API_BASE_URL")

	tpl.ExecuteTemplate(w, "index", data)
}

func RawQuery(w http.ResponseWriter, r *http.Request) {

	var data struct {
		APIBaseURL string
		Query string
	}
	data.APIBaseURL = os.Getenv("API_BASE_URL")
	data.Query = r.URL.RawQuery

	fmt.Println(data)

	tpl.ExecuteTemplate(w, "query", data)
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
