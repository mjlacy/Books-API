package api

import (
	"BookAPI/pkg/database"
	"bookAPI"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
		if err != nil && output.Books != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("An error occurred processing this request")
			return
		}

		if output.Books == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("No books found")
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

		output, err := repo.GetBookById(id)
		if err != nil && output != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("An error occurred processing this request")
			return
		}
		if output == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("No book found with that id")
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
			fmt.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		id, err := repo.PostBook(&u)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "An error occurred processing your request", 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("link: /" + id))
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
			fmt.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		update, err := repo.PutBook(id, &u)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		if update != "" {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("link: /" + update))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("link: /" + id))
		}
	}
}

func Delete(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		err := repo.DeleteBook(id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "not found", 404)
			return
		} else {
			repo.DeleteBook(id)
			w.Write([]byte("Deleted: " + id))
		}
	}
}
