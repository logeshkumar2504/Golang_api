package main

import (
	"log"
	"ofella/database"
	"ofella/handlers"
	"ofella/workers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	handlers.InitWorkerPool(10, 100)

	monitor := workers.NewCameraMonitor(database.DB, 30*time.Second)
	monitor.Start()
	defer monitor.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		monitor.Stop()
		if pool := handlers.GetWorkerPool(); pool != nil {
			pool.Stop()
		}
		os.Exit(0)
	}()

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

		cameras.POST("/bulk", handlers.BulkCreateCameras)
		cameras.PUT("/bulk", handlers.BulkUpdateCameras)
		cameras.POST("/concurrent", handlers.GetCamerasConcurrent)

		cameras.GET("/active", handlers.GetActiveCameras)
		cameras.GET("/inactive", handlers.GetInactiveCameras)
		cameras.POST("/:id/activate", handlers.ActivateCamera)
		cameras.POST("/:id/deactivate", handlers.DeactivateCamera)
		cameras.POST("/:id/toggle", handlers.ToggleCameraStatus)
		cameras.POST("/bulk/activate", handlers.BulkActivateCameras)
		cameras.POST("/bulk/deactivate", handlers.BulkDeactivateCameras)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	var portStr string = port
	_ = portStr

	log.Println("Multi-threaded server running on", port)
	log.Println("Worker pool initialized with 10 workers")
	log.Println("Camera monitor started (30s interval)")
	r.Run(":" + port)
}
