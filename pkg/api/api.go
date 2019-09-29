package api

import (
	"bookAPI"
	"bookAPI/pkg/database"

	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func HealthCheck(repo *database.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.Ping()
		if err != nil {
			log.Println("Error connecting to database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error connecting to the database"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "You have accessed an invalid URL", http.StatusNotFound)
}

func Get(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookId int
		var year int
		var err error
		if bookIdString := r.URL.Query().Get("bookId"); bookIdString != "" {
			bookId, err = strconv.Atoi(bookIdString)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("bookId query must be a nonzero positive integer")
				return
			}
		}

		if yearString := r.URL.Query().Get("year"); yearString != "" {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("year query must be a nonzero positive integer")
				return
			}
		}

		search := bookAPI.Book {
			BookId: bookId,
			Title: r.URL.Query().Get("title"),
			Author: r.URL.Query().Get("author"),
			Year: year,
		}

		output, err := repo.GetBooks(search)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("An error occurred processing this request")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	}
}

func GetById(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		output, err := repo.GetBookById(id)
		if err != nil {
			log.Println(err)
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
	return func(w http.ResponseWriter, r *http.Request) {
		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id, err := repo.PostBook(&u)
		if err != nil {
			log.Println(err)
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}

		u.Id = id

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func Put(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var u = bookAPI.Book{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		update, err := repo.PutBook(id, &u)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if update != "" {
			u.Id = update
		} else {
			u.Id = id
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(u.Id))

		if update != "" {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(u)
	}
}

func Patch(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		//decoder := json.NewDecoder(r.Body)
		//decoder.UseNumber() // will convert to int64 or double

		var update map[string]interface{}

		// err := decoder.Decode(&update)
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = repo.PatchBook(id, update)
		if err != nil {
			log.Println(err)
			if err == bookAPI.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusOK)
	}
}

func Delete(repo bookAPI.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := repo.DeleteBook(id)
		if err != nil {
			log.Println(err)
			if err == bookAPI.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "An error occurred processing your request", http.StatusInternalServerError)
			return
		}
	}
}
