package main

import (
	"BookAPI/internal"

	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockRepository struct {
	id  string
	b   internal.Books
	err error
}

func (r mockRepository) Ping() error {
	return r.err
}

func (r mockRepository) GetBooks(s internal.Book) (b internal.Books, err error) {
	b = r.b
	err = r.err
	return
}

func (r mockRepository) GetBookById(id string) (b *internal.Book, err error) {
	if len(r.b.Books) != 0 {
		b = &r.b.Books[0]
	} else {
		b = nil
	}
	err = r.err
	return
}

func (r mockRepository) PostBook(book *internal.Book) (id string, returnedBook *internal.Book, err error) {
	err = r.err
	return
}

func (r mockRepository) PutBook(id string, book *internal.Book) (bool, *internal.Book, error) {
	if r.b.Books != nil {
		return true, &r.b.Books[0], r.err
	}
	return false, &internal.Book{}, r.err
}

func (r mockRepository) PatchBook(id string, update internal.Book) (err error) {
	if len(r.b.Books) == 0 && r.err == nil {
		return internal.ErrNotFound
	}
	err = r.err
	return
}

func (r mockRepository) DeleteBook(id string) (err error) {
	err = r.err
	return
}

func TestGetBooksSuccess(t *testing.T) {
	r := mockRepository{b: internal.Books{Books: []internal.Book{{}}}}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestGetBooksBadBookId(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodGet, "/?bookId=NaN", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestGetBooksBadYear(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodGet, "/?year=NaN", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestGetBooksError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestGetBooksNoBooksFound(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Get(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestGetBookByIdSuccess(t *testing.T) {
	r := mockRepository{b: internal.Books{Books: []internal.Book{{}}}}

	req, err := http.NewRequest(http.MethodGet, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestGetBookByIdError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest(http.MethodGet, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestGetBookNotFound(t *testing.T) {
	r := mockRepository{err:internal.ErrNotFound}

	req, err := http.NewRequest(http.MethodGet, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := GetById(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func TestCreateBookSuccess(t *testing.T) {
	r := mockRepository{}

	b := internal.Book{BookId: 2, Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestCreateBookBadInput(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestCreateBookError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	b := internal.Book{BookId: 2, Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Post(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestUpsertBookUpdateSuccess(t *testing.T) {
	r := mockRepository{id:"5a80868574fdd6de0f4fa438"}

	b := internal.Book{BookId: 2, Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpsertBookCreateSuccess(t *testing.T) {
	r := mockRepository{b: internal.Books{Books: []internal.Book{{}}}}

	b := internal.Book{BookId: 2, Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpsertBookBadInput(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestUpsertBookError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	b := internal.Book{BookId: 2, Title: "War and Peace", Author: "Leo Tolstoy", Year: 1869}

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPut, "/", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Put(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestUpdateBookSuccess(t *testing.T) {
	r := mockRepository{b: internal.Books{Books: []internal.Book{{}}}}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
	if rr.Header()["Location"] == nil {
		t.Error("No Location header found")
	}
}

func TestUpdateBookBadInput(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodPatch, "/", bytes.NewBuffer([]byte("Bad Input")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %v", rr.Code)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	r := mockRepository{}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5a80868574fdd6de0f4fa438"})

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 but got %v", rr.Code)
	}
}

func TestUpdateBookError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	b := make(map[string]interface{})
	b["bookId"] = 4

	s, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, "/5a80868574fdd6de0f4fa438", bytes.NewBuffer(s))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Patch(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}

func TestDeleteBookSuccess(t *testing.T) {
	r := mockRepository{b: internal.Books{Books: []internal.Book{{}}}}

	req, err := http.NewRequest(http.MethodDelete, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	r := mockRepository{}

	req, err := http.NewRequest(http.MethodDelete, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 but got %v", rr.Code)
	}
}

func TestDeleteBookError(t *testing.T) {
	r := mockRepository{err: errors.New("test error")}

	req, err := http.NewRequest(http.MethodDelete, "/5a80868574fdd6de0f4fa438", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Delete(r)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500 but got %v", rr.Code)
	}
}
