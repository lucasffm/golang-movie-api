package handlers

import (
	"movie-api/config"
	"movie-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateReviewInput struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	Rating    int       `json:"rating" binding:"required,min=1,max=10"`
	Comment   string    `json:"comment"`
	WatchedAt time.Time `json:"watched_at"`
}

type UpdateReviewInput struct {
	Rating    int       `json:"rating" binding:"omitempty,min=1,max=10"`
	Comment   string    `json:"comment"`
	WatchedAt time.Time `json:"watched_at"`
}

func GetReviews(c *gin.Context) {
	var reviews []models.Review
	config.DB.Preload("Movie").Find(&reviews)
	c.JSON(http.StatusOK, reviews)
}

func GetReview(c *gin.Context) {
	var review models.Review
	if err := config.DB.Preload("Movie").First(&review, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}
	c.JSON(http.StatusOK, review)
}

func GetMovieReviews(c *gin.Context) {
	var reviews []models.Review
	if err := config.DB.Where("movie_id = ?", c.Param("id")).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context) {
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var movie models.Movie
	if err := config.DB.First(&movie, input.MovieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	watchedAt := input.WatchedAt
	if watchedAt.IsZero() {
		watchedAt = time.Now()
	}

	review := models.Review{
		MovieID:   input.MovieID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		WatchedAt: watchedAt,
	}

	config.DB.Create(&review)
	config.DB.Preload("Movie").First(&review, review.ID)
	c.JSON(http.StatusCreated, review)
}

func UpdateReview(c *gin.Context) {
	var review models.Review
	if err := config.DB.First(&review, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	var input UpdateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}
	if input.Rating != 0 {
		updateData["rating"] = input.Rating
	}
	if input.Comment != "" {
		updateData["comment"] = input.Comment
	}
	if !input.WatchedAt.IsZero() {
		updateData["watched_at"] = input.WatchedAt
	}

	config.DB.Model(&review).Updates(updateData)
	config.DB.Preload("Movie").First(&review, review.ID)
	c.JSON(http.StatusOK, review)
}

func DeleteReview(c *gin.Context) {
	var review models.Review
	if err := config.DB.First(&review, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	config.DB.Delete(&review)
	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
