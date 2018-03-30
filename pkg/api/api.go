package api

import (
	"net/http"
	"encoding/json"
	"BookAPI/pkg/models"
	"BookAPI/pkg/database"
	"github.com/gorilla/mux"
	"strconv"
)

func HealthCheck(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func NotFoundPage(w http.ResponseWriter, r *http.Request){
	http.Error(w, "You have accessed an invalid URL", 404)
}

func Get(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		output, err := repo.GetBook()
		if err != nil && output.Books == nil {
			http.Error(w, "not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	}
}

func GetById(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		output, err := repo.GetBookById(id)
		if err != nil && output == nil {
			http.Error(w, "not found", 404)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	}
}

func Post(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var u = models.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil{
			http.Error(w, err.Error(), 400)
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

		var u = models.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		idNum, _ := strconv.ParseInt(id, 10, 32) //convert URL id string to int32
		result := int32(idNum)

		test, err := repo.GetBookById(strconv.Itoa(int(u.BookId))) //check if a book exists with the given id in JSON body

		if err != nil && err.Error() != "not found" {
			http.Error(w, "Error checking database for existing record.", 500)
			return
		}

		if err == nil && test != nil && u.BookId != result {
			http.Error(w, "A book already exists with the given Id.", 400)
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
