package routes

import (
	"github.com/gorilla/mux"
	"BookAPI/pkg/api"
	"BookAPI/pkg/database"
	"net/http"
)

type Router struct {
	*mux.Router
}

func New() *Router{
	return &Router{
		mux.NewRouter().StrictSlash(true),
	}
}
func (r *Router)CreateRoutes(db *database.Repository){
	r.HandleFunc("/health", api.HealthCheck).Methods("GET")

	r.HandleFunc("/", api.Get(db)).Methods("GET")

	r.HandleFunc("/{id}", api.GetById(db)).Methods("GET")

	r.HandleFunc("/", api.Post(db)).Methods("POST")

	r.HandleFunc("/{id}", api.Put(db)).Methods("PUT")

	r.HandleFunc("/{id}", api.Delete(db)).Methods("DELETE")

	r.NotFoundHandler = http.HandlerFunc(api.NotFoundPage)
}
