package routes

import (
	"fmt"
	"movie-api/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, movieHandler *handler.MovieHandler, reviewHandler *handler.ReviewHandler) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		fmt.Println("Passou pelo Middleware!!!!")
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		movies := api.Group("/movies")
		{
			movies.GET(
				"/",
				gin.BasicAuth(gin.Accounts{
					"admin": "admin",
				}),
				movieHandler.GetMovies,
			)
			movies.GET("/:id", movieHandler.GetMovie)
			movies.POST("", movieHandler.CreateMovie)
			movies.PUT("/:id", movieHandler.UpdateMovie)
			movies.DELETE("/:id", movieHandler.DeleteMovie)
			movies.GET("/:id/reviews", reviewHandler.GetMovieReviews)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/", reviewHandler.GetReviews)
			reviews.GET("/:id", reviewHandler.GetReview)
			reviews.POST("", reviewHandler.CreateReview)
			reviews.PUT("/:id", reviewHandler.UpdateReview)
			reviews.DELETE("/:id", reviewHandler.DeleteReview)
		}
	}
}
