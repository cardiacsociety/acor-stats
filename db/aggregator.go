package db

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2/bson"
)

// Aggregate runs an aggregation query with the supplied aggregation pipline doc 'ap' and retuns
// a []byte representing the JSON result
func Aggregate(ap []bson.M) ([]bson.M, error) {

	md := os.Getenv("MONGO_DB")
	mc := os.Getenv("MONGO_COL")
	col := mongo.DB(md).C(mc)

	//fmt.Println("Have pipeline doc:", ap)
	pipe := col.Pipe(ap)

	r := []bson.M{}
	err := pipe.All(&r)
	if err != nil {
		fmt.Println("Aggregation error", err)
		return []bson.M{}, err
	}

	return r, nil
}
