package main

import (
	"log"
	"movie-api/config"
	"movie-api/models"
	"movie-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := models.Migrate(config.DB); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	router := gin.Default()
	routes.SetupMiddleware(router)
	routes.SetupRoutes(router)

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
