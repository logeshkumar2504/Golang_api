package main

import (
	"log"
	"ofella/database"
	"ofella/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	cameras := api.Group("/cameras")
	{
		cameras.POST("", handlers.CreateCamera)
		cameras.GET("", handlers.GetAllCameras)
		cameras.GET("/:id", handlers.GetCamera)
		cameras.PUT("/:id", handlers.UpdateCamera)
		cameras.DELETE("/:id", handlers.DeleteCamera)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on", port)
	r.Run(":" + port)
}
