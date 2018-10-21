package bookAPI

import "gopkg.in/mgo.v2/bson"

type Book struct {
	Id                   bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	BookId               int32         `json:"bookId" bson:"bookId"`
	Title                string        `json:"title" bson:"title"`
	Author               string        `json:"author" bson:"author"`
	Year                 int32         `json:"year" bson:"year"`
}

type Books struct {
	Books				[]Book         `json:"Books" bson:"_Books"`
}

type Repository interface {
	GetBooks() (Books, error)
	GetBookByBookId(id string) (b *Book, err error)
	PostBook(book *Book) (err error)
	PutBook(id string, book *Book) (err error)
	DeleteBook(id string) (err error)
}
