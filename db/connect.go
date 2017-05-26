package db

import (
	"fmt"
	"os"

	"github.com/34South/envr"
	"gopkg.in/mgo.v2"
)

var mongo *mgo.Session

func init() {
	// Set up env
	envr.New("acor-db", []string{
		"DEVICES_CSV_FILE",
		"PROCEDURES_CSV_FILE",
		"MONGODB_URI",
		"MONGO_DB",
		"MONGO_COL",
		"MONGO_SRC",
	}).Auto()
}

func Connect() {

	var err error
	mongo, err = mgo.Dial(os.Getenv("MONGODB_URI"))
	if err != nil {
		fmt.Println("Could not connect to mongo", err)
	} else {
		fmt.Println("Connected to mongo", os.Getenv("MONGO_SRC"))
	}
}
