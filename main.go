package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/34South/envr"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
	"time"
	"net/http"
	"io"
	"html/template"
	"github.com/rs/cors"
)

const deviceCSV = "ACOR-Devices.csv"
const proceduresCSV = "ACOR-Procedures.csv"
const reportJSON = "report.json"

type data struct {
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
	PatientID     string    `json:"patientId" bson:"patientId"`
	SiteID        string    `json:"siteId" bson:"siteId"`
	SiteState     string    `json:"siteState" bson:"siteState"`
	ProcDate      time.Time `json:"procDate" bson:"procDate"`
	ProcType      string    `json:"procType" bson:"procType"`
	DeviceType    string    `json:"deviceType" bson:"deviceType"`
	DeviceSubType string    `json:"deviceSubType" bson:"deviceSubType"`
}

type chart struct {
	Title string `json:"title"`
	Scope string `json:"scope"`
	Labels []string `json:"labels"`
	Data []int `json:"data"`
}

var mongo *mgo.Session
var tpl *template.Template
var reportData = []data{}

func init() {

	// Set up env
	envr.New("acor-stats", []string{
		"MONGO_URL",
		"MONGO_DB",
		"MONGO_SRC",
	}).Auto()

	var err error
	mongo, err = mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		fmt.Println("Could not connect to mongo", err)
	} else {
		fmt.Println("Connected to mongo", os.Getenv("MONGO_SRC"))
	}

	tpl = template.Must(template.ParseGlob("./templates/*"))
}

func main() {

	devicesJSON()
	proceduresJSON()
	//JSONFile()
	updateCollection()
	//aggReport("device-NSW")



	mux := http.NewServeMux()

	// API of sorts
	mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("/api/test", testHandler)

	// Pages
	mux.HandleFunc("/report/test", reportHandler)

	// CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
	})
	handler := c.Handler(mux)

	// Specify port when env var is not set - Heroku sets dynamically so cannot include in .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("API listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func devicesJSON() {

	fmt.Println("Looking for csv file", deviceCSV)
	f, err := os.Open(deviceCSV)
	if err != nil {
		log.Fatalln("Could not open", deviceCSV, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		xs := strings.Split(s.Text(), ",")
		d := data{}
		d.UpdatedAt = time.Now()
		d.PatientID = xs[0]
		d.SiteID = xs[1]
		d.SiteState = xs[2]
		t, err := time.Parse("2/01/2006", xs[3])
		if err != nil {
			fmt.Println("time.Parse() err,", err)
		}
		d.ProcDate = t
		d.ProcType = "device"
		if xs[5] == "" {
			d.DeviceType = "Unknown"
		} else {
			d.DeviceType = xs[5]
		}
		if xs[7] == "" {
			d.DeviceSubType = "Unknown"
		} else {
			d.DeviceSubType = xs[7]
		}

		reportData = append(reportData, d)
	}
}

func proceduresJSON() {

	fmt.Println("Looking for csv file", proceduresCSV)
	f, err := os.Open(proceduresCSV)
	if err != nil {
		log.Fatalln("Could not open", proceduresCSV, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		xs := strings.Split(s.Text(), ",")
		d := data{}
		d.UpdatedAt = time.Now()
		d.PatientID = xs[0]
		d.SiteID = xs[1]
		d.SiteState = xs[2]
		t, err := time.Parse("2/01/2006", xs[3])
		if err != nil {
			fmt.Println("time.Parse() err,", err)
		}
		d.ProcDate = t
		d.ProcType = "pci"
		d.DeviceType = "stent"
		if xs[6] == "" {
			d.DeviceSubType = "Unknown"
		} else {
			d.DeviceSubType = xs[6]
		}

		reportData = append(reportData, d)
	}
}

func JSONFile() {

	//j, err := json.MarshalIndent(reportData, "", "  ")
	//if err != nil {
	//	fmt.Println("Could not marshal data", err)
	//	os.Exit(1)
	//}

	//fmt.Println(string(j))
	//err = ioutil.WriteFile(reportJSON, j, 0644)
	//if err != nil {
	//	fmt.Println("Could not create file", reportJSON)
	//	os.Exit(1)
	//}
}

func updateCollection() {

	data := mongo.DB("acor").C("Data")

	s := map[string]interface{}{}
	data.Find(bson.M{}).One(&s)

	for _, v := range reportData {
		//fmt.Println("Upsert", i, v)

		// Use patientId, siteId and procDate as selector for upsert
		s := bson.M{"patientId": v.PatientID, "siteId": v.SiteID, "procDate": v.ProcDate}

		_, err := data.Upsert(s, v)
		if err != nil {
			fmt.Println("Error inserting doc", err)
		}
	}
}

func aggReport(a string) []byte {

	data := mongo.DB("acor").C("Data")

	// get aggregation...
	ap := aggPipe(a)
	fmt.Println("Have pipeline doc:", ap)
	pipe := data.Pipe(ap)

	r := []bson.M{}
	err := pipe.All(&r)
	if err != nil {
		fmt.Println("Aggregation error", err)
		return []byte{}
	}

	// JSON the report
	//j, err := json.MarshalIndent(r, "", "  ")
	//j, err := json.MarshalIndent(r, "", "  ")
	//if err != nil {
	//	fmt.Println("Could not marshal json report", err)
	//}
	//fmt.Println(string(j))


	// For reports ideally we just need a JSON with label / value pairs, in order
	rj := chart{Title: "Device for NSW", Scope: "Devices, NSW"}
	//months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Nov", "Dec"}
	for i, v := range r {
		point := v["_id"].(bson.M)
		fmt.Println(i, point["month"], "-", point["year"])
		//monthNumber := point["month"].(int)
		rj.Labels = append(rj.Labels, fmt.Sprintf("%v-%v", point["month"], point["year"]))
		rj.Data = append(rj.Data, v["count"].(int))
	}

	//fmt.Println(rj)

	xb, err := json.Marshal(rj)
	if err != nil {
		fmt.Println("Could not marshal json report", err)
	}
	//fmt.Println(string(xb))

	return xb

}

func aggPipe(a string) []bson.M {

	fmt.Println("Generating report:", a)

	switch a {

	case "device-NSW":
		return []bson.M{
			{"$match": bson.M{"siteState": "NSW", "procType": "device"}},
			{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
			{"$sort": bson.M{"_id.year": 1, "_id.month": 1}},
		}

	case "device-QLD":
		return []bson.M{
			{"$match": bson.M{"siteState": "QLD", "procType": "device"}},
			{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
		}

	case "device-SA":
		return []bson.M{
			{"$match": bson.M{"siteState": "SA", "procType": "device"}},
			{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
		}

	case "device-WA":
		return []bson.M{
			{"$match": bson.M{"siteState": "WA", "procType": "device"}},
			{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
		}

	}

	// return to match nothing
	return []bson.M{
		{"$match": bson.M{"siteState": "XXXXXX", "procType": "XXXXX"}},
		{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
	}
}



func testHandler(w http.ResponseWriter, r *http.Request) {

	xb := aggReport("device-NSW")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, string(xb))
}

func reportHandler(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "test", nil)
}