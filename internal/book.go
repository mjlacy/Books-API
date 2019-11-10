package internal

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// ErrNotFound indicates that a book that matches the given criteria was not found
	ErrNotFound = errors.New("Book not found")
	// ErrInvalidId indicates an invalid id was sent to the application
	ErrInvalidId = errors.New("Invalid id provided")
)

// Repository defines a basic set of actions to interact with books at the data source
type Repository interface {
	Ping() error
	GetBooks(book Book) (books Books, err error)
	GetBookById(id string) (book *Book, err error)
	PostBook(book *Book) (id string, returnedBook *Book, err error)
	PutBook(id string, book *Book) (isCreated bool, returnedBook *Book, err error)
	PatchBook(id string, update Book) (err error)
	DeleteBook(id string) (err error)
}

// Books models a list of books displayed together
type Books struct {
	Books []Book `json:"books"`
}

// Book models the structure of a book stored in the data source
type Book struct {
	Id     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	BookId int32              `json:"bookId" bson:"bookId"`
	Title  string             `json:"title" bson:"title"`
	Author string             `json:"author" bson:"author"`
	Year   int32              `json:"year" bson:"year"`
}
