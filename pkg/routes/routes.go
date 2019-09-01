package routes

import (
	"bookAPI/pkg/api"
	"bookAPI/pkg/database"

	"net/http"

	"github.com/gorilla/mux"
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
	r.HandleFunc("/health", api.HealthCheck(db)).Methods(http.MethodGet)

	r.HandleFunc("/", api.Get(db)).Methods(http.MethodGet)

	r.HandleFunc("/{id}", api.GetById(db)).Methods(http.MethodGet)

	r.HandleFunc("/", api.Post(db)).Methods(http.MethodPost)

	r.HandleFunc("/{id}", api.Put(db)).Methods(http.MethodPut)

	r.HandleFunc("/{id}", api.Patch(db)).Methods(http.MethodPatch)

	r.HandleFunc("/{id}", api.Delete(db)).Methods(http.MethodDelete)

	r.NotFoundHandler = http.HandlerFunc(api.NotFoundPage)
}
