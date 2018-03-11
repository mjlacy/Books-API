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

	http.ListenAndServe(configs.ThisPortNumber, cors.Default().Handler(r))
}
