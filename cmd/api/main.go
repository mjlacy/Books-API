package main

import (
	"BookAPI/internal/mongo"
	"BookAPI/pkg/configuration"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	configs, err := configuration.New("cmd/api/config.json")
	if err != nil {
		log.Fatalln("Error opening properties file: ", err)
	}

	db, err := mongo.InitializeMongoDatabase(&mongo.DatabaseConfig{
		DbURL: configs.DbURL,
		DatabaseName: configs.DatabaseName,
		CollectionName: configs.CollectionName})
	if err != nil {
		log.Fatalln("Error connecting to mongo: ", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheck(db)).Methods(http.MethodGet)
	r.HandleFunc("/", Get(db)).Methods(http.MethodGet)
	r.HandleFunc("/{id}", GetById(db)).Methods(http.MethodGet)
	r.HandleFunc("/", Post(db)).Methods(http.MethodPost)
	r.HandleFunc("/{id}", Put(db)).Methods(http.MethodPut)
	r.HandleFunc("/{id}", Patch(db)).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", Delete(db)).Methods(http.MethodDelete)
	r.NotFoundHandler = http.HandlerFunc(NotFoundPage)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
	handler := cors.Default().Handler(r)
	handler = c.Handler(handler)

	if err := http.ListenAndServe(configs.ThisPortNumber, handler); err != nil {
		log.Fatalln("Error launching api:", err)
	}
}
