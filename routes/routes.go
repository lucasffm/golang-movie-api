package routes

import (
	"fmt"
	"movie-api/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) *gin.Engine {
	// Setup middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		fmt.Println("Passou pelo Middleware!!!!")
		c.Next()
	})

	// Setup routes
	api := router.Group("/api/v1")
	{
		movies := api.Group("/movies")
		{
			movies.GET(
				"/",
				gin.BasicAuth(gin.Accounts{
					"admin": "admin",
				}),
				handlers.GetMovies,
			)
			movies.GET("/:id", handlers.GetMovie)
			movies.POST("", handlers.CreateMovie)
			movies.PUT("/:id", handlers.UpdateMovie)
			movies.DELETE("/:id", handlers.DeleteMovie)
			movies.GET("/:id/reviews", handlers.GetMovieReviews)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("/", handlers.GetReviews)
			reviews.GET("/:id", handlers.GetReview)
			reviews.POST("", handlers.CreateReview)
			reviews.PUT("/:id", handlers.UpdateReview)
			reviews.DELETE("/:id", handlers.DeleteReview)
		}
	}

	return router
}
