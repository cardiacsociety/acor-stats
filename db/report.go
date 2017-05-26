package db

import (
	"gopkg.in/mgo.v2/bson"
)

type ReportFilter struct {
	State         string // the state
	ProcType      string // the type of procedure - "device" for all device registry, or one of the other procedures
	DeviceType    string // the type of device involved
	DeviceSubType string // the sub type of device
}

// Stores all of the report queries / aggregation docs

// TimeReport does reports by month-year. It takes a report filter to modify the aggregation query
func TimeReport(rf ReportFilter) []bson.M {

	//fmt.Println("\nTimeReport() --------------------------------")
	match := bson.M{}
	if rf.State != "" {
		match["siteState"] = rf.State
	}
	if rf.ProcType != "" {
		match["procType"] = rf.ProcType
	}
	if rf.DeviceType != "" {
		match["deviceType"] = rf.DeviceType
	}

	if rf.DeviceSubType != "" {
		match["deviceSubType"] = rf.DeviceSubType
	}

	return []bson.M{
		{"$match": match},
		{"$group": bson.M{"_id": bson.M{"month": bson.M{"$month": "$procDate"}, "year": bson.M{"$year": "$procDate"}}, "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"_id.year": 1, "_id.month": 1}},
	}

}
