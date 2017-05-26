package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cardiacsociety/acor-stats/db"
	"github.com/cardiacsociety/acor-stats/handlers"
	"github.com/rs/cors"
)

const ReportJSON = "report.json"

func main() {

	db.Connect()
	//db.Import()

	//aggReport("device-NSW")

	// Kick off the web server and api
	m := http.NewServeMux()

	// Static
	m.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./public/css"))))
	m.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./public/js"))))
	m.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./public/img"))))
	m.Handle("/favicon.ico", http.NotFoundHandler())

	// API
	m.HandleFunc("/api/report/time", handlers.TimeReportHandler)

	// Pages / Reports
	m.HandleFunc("/", handlers.Index)
	m.HandleFunc("/test", handlers.Test)
	m.HandleFunc("/reports/all", handlers.ReportsAllHandler)
	m.HandleFunc("/reports/all/state", handlers.ReportsAllStateHandler)
	m.HandleFunc("/reports/devices", handlers.ReportsDevicesHandler)
	m.HandleFunc("/reports/devices/state", handlers.ReportsDevicesStateHandler)
	m.HandleFunc("/reports/procedures", handlers.ReportsProceduresHandler)
	m.HandleFunc("/reports/procedures/state", handlers.ReportsProceduresStateHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
	})
	handler := c.Handler(m)

	// Specify port when env var is not set - Heroku sets dynamically so cannot include in .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("API listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
