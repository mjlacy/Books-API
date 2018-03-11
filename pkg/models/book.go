package models

import "gopkg.in/mgo.v2/bson"

type Book struct {
	Id                   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Name                 string        `json:"Name" bson:"Name"`
	Author               string        `json:"Author" bson:"Author"`
	Year                 int32         `json:"Year" bson:"Year"`
}