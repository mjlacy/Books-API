package api

import (
	"bookAPI"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRepository struct {
	id string
	b bookAPI.Books
	err error
}

func (r mockRepository) GetBooks(s bookAPI.Book) (b bookAPI.Books, err error){
	b = r.b
	err = r.err
	return
}

func (r mockRepository) GetBookById(id string) (b *bookAPI.Book, err error){
	if len(r.b.Books) != 0 {
		mId, _ := primitive.ObjectIDFromHex(id)
		if r.b.Books[0].Id == mId {
			b = &r.b.Books[0]
		} else {
			b = nil
		}
	} else {
		b = nil
	}
	err = r.err
	return
}

func (r mockRepository) PostBook(book *bookAPI.Book) (id string, err error){
	err = r.err
	return
}

func (r mockRepository) PutBook(id string, book *bookAPI.Book) (bool, *bookAPI.Book, error){
	//updateId = r.id
	//err := r.err
	if r.b.Books != nil {
		return true, &r.b.Books[0], r.err
	}
	return false, &bookAPI.Book{}, r.err
}

func (r mockRepository) PatchBook(id string, update map[string]interface{}) (err error){
	if len(r.b.Books) != 0 {
		mId, _ := primitive.ObjectIDFromHex(id)
		if r.b.Books[0].Id != mId {
			err = errors.New("not found")
			return
		}
	}
	err = r.err
	return
}

func (r mockRepository) DeleteBook(id string) (err error){
	err = r.err
	return
}

func TestGetBooksSuccess(t *testing.T){
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{}}}}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestGetBooksBadBookId(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("GET", "/?bookId=NaN", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest{
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestGetBooksBadYear(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("GET", "/?year=NaN", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest{
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestGetBooksError(t *testing.T){
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestGetBooksNoBooksFound(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound{
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func TestGetBookByIdSuccess(t *testing.T){
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	req, err := http.NewRequest("GET", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestGetBookByIdError(t *testing.T){
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest("GET", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestGetBookNotFound(t *testing.T){
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa437")
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	req, err := http.NewRequest("GET", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound{
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func TestCreateBookSuccess(t *testing.T){
	r := mockRepository{}

	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	b := bookAPI.Book{Id:mId, BookId: 2,
	  Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated{
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestCreateBookBadInput(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest{
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestCreateBookError(t *testing.T){
	r := mockRepository{err: errors.New("test error")}

	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	b := bookAPI.Book{Id:mId, BookId: 2,
		Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestUpsertBookUpdateSuccess(t *testing.T){
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")

	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	b := bookAPI.Book{Id:mId, BookId: 2,
		Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpsertBookCreateSuccess(t *testing.T){
	r := mockRepository{id:"5a80868574fdd6de0f4fa438"}

	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	b := bookAPI.Book{Id:mId, BookId: 2,
		Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated{
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpsertBookBadInput(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("PUT", "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest{
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestUpsertBookError(t *testing.T){
	r := mockRepository{err: errors.New("test error")}

	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	b := bookAPI.Book{Id:mId, BookId: 2,
		Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestUpdateBookSuccess(t *testing.T){
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpdateBookBadInput(t *testing.T){
	r := mockRepository{}

	req, err := http.NewRequest("PATCH", "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest{
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestUpdateBookNotFound(t *testing.T){
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa437")
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound{
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func TestUpdateBookError(t *testing.T){
	r := mockRepository{err: errors.New("test error")}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestDeleteBookSuccess(t *testing.T) {
	mId, _ := primitive.ObjectIDFromHex("5a80868574fdd6de0f4fa438")
	r := mockRepository{b: bookAPI.Books{Books: []bookAPI.Book{{Id: mId}}}}

	req, err := http.NewRequest("DELETE", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest("DELETE", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestDeleteBookError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest("DELETE", "/5a80868574fdd6de0f4fa438", nil)
	if err != nil{
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError{
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}
