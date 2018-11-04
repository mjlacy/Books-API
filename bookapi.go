package bookAPI

import "gopkg.in/mgo.v2/bson"

type Book struct {
	Id                   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	BookId               int           `json:"bookId" bson:"bookId"`
	Title                string        `json:"title" bson:"title"`
	Author               string        `json:"author" bson:"author"`
	Year                 int           `json:"year" bson:"year"`
}

type Books struct {
	Books				[]Book         `json:"books"`
}

type Repository interface {
	GetBooks(s Book) (b Books, err error)
	GetBookById(id string) (b *Book, err error)
	PostBook(book *Book) (id string, err error)
	PutBook(id string, book *Book) (updateId string, err error)
	DeleteBook(id string) (err error)
}
