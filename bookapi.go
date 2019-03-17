package bookAPI

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	BookId int32             `json:"bookId" bson:"bookId"`
	Title  string            `json:"title" bson:"title"`
	Author string            `json:"author" bson:"author"`
	Year   int32             `json:"year" bson:"year"`
}

type Books struct {
	Books []Book `json:"books"`
}

type Repository interface {
	GetBooks(s Book) (b Books, err error)
	GetBookById(id string) (b *Book, err error)
	PostBook(book *Book) (id string, err error)
	PutBook(id string, book *Book) (bool, *Book, error)
	//PatchBook(id string, update map[string]interface{}) (err error)
	DeleteBook(id string) (err error)
}
