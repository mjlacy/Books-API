package main

import (
	"BookAPI/internal"
	"BookAPI/internal/mongo"

	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func HealthCheck(repo *mongo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := repo.Ping()
		if err != nil {
			log.Printf("Error connecting to database: %v", err)
			http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func NotFoundPage(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "You have accessed an invalid URL", http.StatusNotFound)
}

func Get(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var queryStrings internal.Book
		err := schema.NewDecoder().Decode(&queryStrings, r.URL.Query())
		if err != nil {
			log.Printf("Unable to decode query string: %v", err)
			http.Error(w, "Invalid query string", http.StatusBadRequest)
			return
		}

		books, err := repo.GetBooks(queryStrings)
		if err != nil {
			log.Printf("Error getting books: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(books); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func GetById(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		output, err := repo.GetBookById(id)
		if err != nil {
			if err == internal.ErrNotFound {
				http.Error(w, "No book found with that id", http.StatusNotFound)
				return
			}

			log.Printf("Error getting book: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func Post(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book = internal.Book{}

		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			log.Printf("Unable to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		id, resp, err := repo.PostBook(&book)
		if err != nil {
			log.Printf("Error creating book: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func Put(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var book = internal.Book{}
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			log.Printf("Unable to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		isUpserted, updatedBook, err := repo.PutBook(id, &book)
		if err != nil {
			if err == internal.ErrInvalidId {
				http.Error(w, "The given id is not a valid id", http.StatusBadRequest)
				return
			}

			log.Printf("Error upserting language: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))

		if isUpserted {
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(updatedBook); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}

func Patch(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var update internal.Book

		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Printf("Unable to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = repo.PatchBook(id, update)
		if err != nil {
			if err == internal.ErrNotFound {
				http.Error(w, "No book with that id found to update", http.StatusNotFound)
				return
			}
			log.Printf("Error updating book: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Location", "/" + url.PathEscape(id))
		w.WriteHeader(http.StatusOK)
	}
}

func Delete(repo internal.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := repo.DeleteBook(id)
		if err != nil {
			if err == internal.ErrNotFound {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			log.Printf("Error deleting book: %v", err)
			http.Error(w, "An error occurred processing this request", http.StatusInternalServerError)
			return
		}
	}
}
