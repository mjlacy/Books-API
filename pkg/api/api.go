package api

import (
	"BookAPI/pkg/database"
	"bookAPI"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

func HealthCheck(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		//err := repo.Ping()
		//if err != nil {
		//	fmt.Println("Error connecting to database: ", err)
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte("Error connecting to the database"))
		//}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func NotFoundPage(w http.ResponseWriter, r *http.Request){
	http.Error(w, "You have accessed an invalid URL", http.StatusNotFound)
}

func Get(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var bookId int
		var year int
		var err error
		if bookIdString := r.URL.Query().Get("bookId"); bookIdString != ""{
			bookId, err = strconv.Atoi(bookIdString)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("bookId query must be a nonzero positive integer")
				return
			}
		}

		if yearString := r.URL.Query().Get("year"); yearString != ""{
			year, err = strconv.Atoi(yearString)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("year query must be a nonzero positive integer")
				return
			}
		}

		search := bookAPI.Book{
			BookId: int32(bookId),
			Title: r.URL.Query().Get("title"),
			Author: r.URL.Query().Get("author"),
			Year: int32(year),
		}

		output, err := repo.GetBooks(search)
		if err != nil{
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("An error occurred processing this request")
			return
		}

		//if output.Books == nil {
		if len(output) == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("No books found")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	}
}

func GetById(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		output, err := repo.GetBookById(id)
		if err != nil {
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
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	}
}

func Post(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil{
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := repo.PostBook(&u)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func Put(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		update, err := repo.PutBook(id, &u)
		if err != nil {
			fmt.Println(err)
			if err.Error() == "Invalid id given" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))

		if update != "" {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(u)
	}
}

//func Patch(repo bookAPI.Repository) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		id := mux.Vars(r)["id"]
//
//		//decoder := json.NewDecoder(r.Body)
//		//decoder.UseNumber() // will convert to int64 or double
//
//		var update map[string]interface{}
//
//		// err := decoder.Decode(&update)
//		err := json.NewDecoder(r.Body).Decode(&update)
//		if err != nil {
//			fmt.Println(err)
//			http.Error(w, err.Error(), 400)
//			return
//		}
//
//		err = repo.PatchBook(id, update)
//		if err != nil {
//			fmt.Println(err)
//			if err.Error() == "not found" {
//				http.Error(w, "No book with that _id found to update", 404)
//				return
//			} else if err.Error() == "Invalid id given" {
//				http.Error(w, err.Error(), 400)
//				return
//			}
//			http.Error(w, "An error occurred processing your request", 500)
//			return
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.Header().Add("Location", "/" + url.PathEscape(id))
//		w.WriteHeader(http.StatusOK)
//	}
//}

func Delete(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		id := mux.Vars(r)["id"]

		err := repo.DeleteBook(id)
		if err != nil {
			if err.Error() == "not found" {
				http.Error(w, "not found", http.StatusNotFound)
			} else {
				http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			}
		}
	}
}
