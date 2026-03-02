package main

import (
	"movie-api/config"
	"movie-api/internal/domain"
	"movie-api/internal/handler"
	"movie-api/internal/repository"
	"movie-api/internal/usecase"
	"movie-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Msg("No .env file found, using environment variables")
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Err(err).Msg("Failed to connect to database")
	}

	if err := domain.Migrate(db); err != nil {
		log.Err(err).Msg("Failed to migrate database")
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	movieRepo := repository.NewMovieRepository(db, &logger)
	reviewRepo := repository.NewReviewRepository(db, &logger)

	movieUC := usecase.NewMovieUseCase(movieRepo)
	reviewUC := usecase.NewReviewUseCase(reviewRepo, movieRepo)

	movieHandler := handler.NewMovieHandler(movieUC)
	reviewHandler := handler.NewReviewHandler(reviewUC)

	router := gin.Default()
	routes.Setup(router, movieHandler, reviewHandler)

	log.Info().Msg("Server running on port 8080!!!!")
	if err := router.Run(":8080"); err != nil {
		log.Err(err).Msg("Failed to start server")
	}
}
