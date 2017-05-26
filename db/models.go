package db

import "time"

type Data struct {
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
	PatientID     string    `json:"patientId" bson:"patientId"`
	SiteID        string    `json:"siteId" bson:"siteId"`
	SiteState     string    `json:"siteState" bson:"siteState"`
	ProcDate      time.Time `json:"procDate" bson:"procDate"`
	ProcType      string    `json:"procType" bson:"procType"`
	DeviceType    string    `json:"deviceType" bson:"deviceType"`
	DeviceSubType string    `json:"deviceSubType" bson:"deviceSubType"`
}
