package api

import (
	"BookAPI/pkg/database"
	"bookAPI"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func HealthCheck(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		err := repo.Ping()
		if err != nil {
			fmt.Println("Error connecting to database: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error connecting to the database"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func NotFoundPage(w http.ResponseWriter, r *http.Request){
	http.Error(w, "You have accessed an invalid URL", 404)
}

func Get(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		output, err := repo.GetBooks()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("An error occurred processing this request"))
			return
		}

		if output.Books == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No books found"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(output)
	}
}

func GetById(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		output, err := repo.GetBookByBookId(id)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("An error occurred processing this request"))
			return
		}

		if output == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No book found with that id"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(output)
	}
}

func Post(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil{
			http.Error(w, err.Error(), 400)
			return
		}

		test, err := repo.GetBookByBookId(strconv.Itoa(int(u.BookId))) //check if a book exists with the given id in JSON body

		if err != nil && err.Error() != "not found" {
			http.Error(w, "Error checking database for existing record.", 500)
			return
		}

		if err == nil && test != nil {
			http.Error(w, "A book already exists with the given BookId.", 400)
			return
		}

		err = repo.PostBook(&u)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Put(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		if id == "" {
			http.Error(w, "bad request", 400)
			return
		}

		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		idNum, _ := strconv.ParseInt(id, 10, 32) //convert URL id string to int32
		result := int32(idNum)

		test, err := repo.GetBookByBookId(strconv.Itoa(int(u.BookId))) //check if a book exists with the given id in JSON body

		if err != nil && err.Error() != "not found" {
			http.Error(w, "Error checking database for existing record.", 500)
			return
		}

		if err == nil && test != nil && u.BookId != result {
			http.Error(w, "A book already exists with the given BookId.", 400)
			return
		}

		err = repo.PutBook(id, &u)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Delete(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		err := repo.DeleteBook(id)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		} else{
			repo.DeleteBook(id)
			w.Write([]byte("Deleted: " + id))
		}
	}
}
