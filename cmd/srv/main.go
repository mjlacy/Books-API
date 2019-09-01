package main

import (
	"bookAPI/pkg/configuration"
	"bookAPI/pkg/database"
	"bookAPI/pkg/routes"

	"log"
	"net/http"

	"github.com/rs/cors"
)

func main(){
	configs, err := configuration.New("cmd/srv/config.json") // cmd/srv/config.json in IntelliJ, config.json in VSCode
	if err != nil {
		log.Println("Error opening properties file: ", err)
	}

	db, err := database.InitializeMongoDatabase(&database.DatabaseConfig{
		DbURL:          configs.DbURL,
		DatabaseName:   configs.DatabaseName,
		CollectionName: configs.CollectionName})
	if err != nil {
		log.Println("Error connecting to database: ", err)
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
