package main

import (
	"BookAPI/pkg/routes"
	"BookAPI/pkg/configuration"
	"BookAPI/pkg/database"
	"github.com/rs/cors"
	"net/http"
)

func main(){
	configs := configuration.New("cmd/srv/config.json")

	db := database.InitializeMongoDatabase(&database.DatabaseConfig{
		DbURL: configs.DbURL,
		DatabaseName: configs.DatabaseName,
		CollectionName: configs.CollectionName,
		})

	r := routes.New()
	r.CreateRoutes(db)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := cors.Default().Handler(r)
	handler = c.Handler(handler);

	http.ListenAndServe(configs.ThisPortNumber, handler)
}
