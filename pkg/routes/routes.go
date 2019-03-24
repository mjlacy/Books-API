package routes

import (
	"BookAPI/pkg/api"
	"BookAPI/pkg/database"
	"github.com/gin-gonic/gin"
)

func New(db *database.Repository) *gin.Engine{
	router := gin.New()

	//router.GET("/health", api.HealthCheck(db))

	router.GET("/", api.Get(db))

	router.GET("/:id", api.GetById(db))

	router.POST("/", api.Post(db))

	router.PUT("/:id", api.Put(db))

	router.PATCH("/:id", api.Patch(db))

	router.DELETE("/:id", api.Delete(db))

	router.NoRoute()

	return router
}
