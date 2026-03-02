package main

import (
	"log"
	"movie-api/config"
	"movie-api/internal/domain"
	"movie-api/internal/handler"
	"movie-api/internal/repository"
	"movie-api/internal/usecase"
	"movie-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := domain.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	movieRepo := repository.NewMovieRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	movieUC := usecase.NewMovieUseCase(movieRepo)
	reviewUC := usecase.NewReviewUseCase(reviewRepo, movieRepo)

	movieHandler := handler.NewMovieHandler(movieUC)
	reviewHandler := handler.NewReviewHandler(reviewUC)

	router := gin.Default()
	routes.Setup(router, movieHandler, reviewHandler)

	log.Println("Server running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
