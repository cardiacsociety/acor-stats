package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Import will parse the CSV files and upsert data to the mongo collection
func Import() {

	reportData := []Data{}
	f := os.Getenv("DEVICES_CSV_FILE")
	importDevicesCSV(&reportData, f)
	f = os.Getenv("PROCEDURES_CSV_FILE")
	importProceduresCSV(&reportData, f)

	// Write the data to a JSON file
	//JSONFile()

	// Upsert to mongo db
	updateCollection(reportData)

}

func importDevicesCSV(rd *[]Data, csv string) {

	fmt.Println("Looking for csv file", csv)
	f, err := os.Open(csv)
	if err != nil {
		log.Fatalln("Could not open", csv, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		xs := strings.Split(s.Text(), ",")
		d := Data{}
		d.UpdatedAt = time.Now()
		d.SiteState = xs[0]
		t, err := time.Parse("2/01/2006", xs[1])
		if err != nil {
			fmt.Println("time.Parse() err,", err)
		}
		d.ProcDate = t
		d.ProcType = "device"
		d.DeviceType = xs[3]
		if d.DeviceType == "" {
			d.DeviceType = "Unknown"
		}
		d.DeviceSubType = xs[5]
		if d.DeviceSubType == "" {
			d.DeviceSubType = "Unknown"
		}

		*rd = append(*rd, d)
	}
}

func importProceduresCSV(rd *[]Data, csv string) {

	fmt.Println("Looking for csv file", csv)
	f, err := os.Open(csv)
	if err != nil {
		log.Fatalln("Could not open", csv, err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		xs := strings.Split(s.Text(), ",")
		d := Data{}
		d.UpdatedAt = time.Now()
		d.SiteState = xs[0]
		t, err := time.Parse("2/01/2006", xs[1])
		if err != nil {
			fmt.Println("time.Parse() err,", err)
		}
		d.ProcDate = t
		d.ProcType = "pci"
		d.DeviceType = "stent"
		d.DeviceSubType = xs[3]
		if d.DeviceSubType == "" {
			d.DeviceSubType = "None"
		}

		*rd = append(*rd, d)
	}
}

// ExportJSON will write data to a JSOn file
//func ExportJSON() {

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
//}

func updateCollection(d []Data) {

	md := os.Getenv("MONGO_DB")
	mc := os.Getenv("MONGO_COL")
	col := mongo.DB(md).C(mc)

	//s := map[string]interface{}{}
	//col.Find(bson.M{}).One(&s)
	//fmt.Println(s)

	// Straight insert means we need to empty out first
	_, err := col.RemoveAll(bson.M{})
	if err != nil {
		log.Fatalln("Could not empty collection")

	}

	for _, v := range d {
		//fmt.Println("Upsert", i, v)

		// Use patientId, siteId and procDate as selector for upsert
		//s := bson.M{"patientId": v.PatientID, "siteId": v.SiteID, "procDate": v.ProcDate}
		//_, err := col.Upsert(s, v)
		//if err != nil {
		//	fmt.Println("Error inserting doc", err)
		//}

		if err := col.Insert(v); err != nil {
			fmt.Println("Error inserting doc", err)
		}
		//fmt.Println(s)
	}
}
