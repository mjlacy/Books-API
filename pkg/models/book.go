package models

//import "gopkg.in/mgo.v2/bson"

type Book struct {
	Id                   string        `json:"_id" bson:"_id,omitempty"`
	BookId               int32         `json:"BookId" bson:"BookId"`
	Name                 string        `json:"Name" bson:"Name"`
	Author               string        `json:"Author" bson:"Author"`
	Year                 int32         `json:"Year" bson:"Year"`
}

type Books struct {
	Books				[]Book         `json:"Books" bson:"_Books"`
}