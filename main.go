package main

import (
	"ai-incident-manager/database"
	"ai-incident-manager/handlers"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	database.Init()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: false,
	}))

	r.StaticFile("/", "./static/index.html")

	api := r.Group("/api")
	{
		api.GET("/incidents", handlers.ListIncidents)
		api.GET("/incidents/:id", handlers.GetIncident)
		api.POST("/incidents", handlers.CreateIncident)
		api.PUT("/incidents/:id", handlers.UpdateIncident)
		api.DELETE("/incidents/:id", handlers.DeleteIncident)
		api.POST("/incidents/:id/analyze", handlers.AnalyzeIncident)
		api.GET("/incidents/:id/messages", handlers.GetMessages)
		api.POST("/incidents/:id/chat", handlers.Chat)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
