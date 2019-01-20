package main

import (
	"BookAPI/pkg/configuration"
	"BookAPI/pkg/database"
	"BookAPI/pkg/routes"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func main(){
	configs, err := configuration.New("cmd/srv/config.json")
	if err != nil {
		fmt.Println("Error opening properties file: ", err)
	}

	db, err := database.InitializeMongoDatabase(&database.DatabaseConfig{
		DbURL: configs.DbURL,
		DatabaseName: configs.DatabaseName,
		CollectionName: configs.CollectionName})
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
	}

	r := routes.New()
	r.CreateRoutes(db)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	})
	handler := cors.Default().Handler(r)
	handler = c.Handler(handler)

	http.ListenAndServe(configs.ThisPortNumber, handler)
}
