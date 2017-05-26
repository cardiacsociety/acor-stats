package handlers

// Provides the handlers for API calls - all responses are JSON

import (
	"io"
	"net/http"

	"encoding/json"
	"fmt"
	"github.com/cardiacsociety/acor-stats/db"
	"gopkg.in/mgo.v2/bson"
)

type Chart struct {
	Title  string   `json:"title"`
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

func TimeReportHandler(w http.ResponseWriter, r *http.Request) {

	// get state param
	//report := r.URL.Query().Get("report")

	fmt.Println(r.URL.String())
	fmt.Println(r.URL.RawQuery)

	// filter object
	rf := db.ReportFilter{
		State:         r.URL.Query().Get("s"),
		ProcType:      r.URL.Query().Get("p"),
		DeviceType:    r.URL.Query().Get("d"),
		DeviceSubType: r.URL.Query().Get("ds"),
	}

	// Make the aggregation pipeline based on the filter
	ap := db.TimeReport(rf)
	//fmt.Println(ap)

	// Run the query, returns results as []bson.M
	xbm, err := db.Aggregate(ap)
	if err != nil {
		RespondError(w, err)
		return
	}

	// Pull out the bits we need and turn it into a JSON string
	title := fmt.Sprintf("Test Report: - %v - %v - %v - %v", rf.State, rf.ProcType, rf.DeviceType, rf.DeviceSubType)
	xb, err := chartJSON(title, xbm)
	if err != nil {
		RespondError(w, err)
		return
	}

	Respond(w, string(xb))
}

// chartJSON returns chart data json as []byte
func chartJSON(title string, xbm []bson.M) ([]byte, error) {

	// Add title...
	rj := Chart{Title: title}
	//months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Nov", "Dec"}

	// ... and array of labels and values, in order
	for i, v := range xbm {
		point := v["_id"].(bson.M)
		fmt.Println(i, point["month"], "-", point["year"])
		//monthNumber := point["month"].(int)
		rj.Labels = append(rj.Labels, fmt.Sprintf("%v-%v", point["month"], point["year"]))
		rj.Data = append(rj.Data, v["count"].(int))
	}

	xb, err := json.Marshal(rj)
	if err != nil {
		fmt.Println("Could not marshal json report", err)
	}
	//fmt.Println(string(xb))

	return xb, err
}

func Respond(w http.ResponseWriter, body string) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(w, body)
}

func RespondError(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body := fmt.Sprintf(`{"error": "%v"}`, err)
	io.WriteString(w, body)
}
