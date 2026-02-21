package handlers

import (
	"movie-api/config"
	"movie-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateMovieInput struct {
	Title       string `json:"title" binding:"required"`
	Director    string `json:"director"`
	Year        int    `json:"year"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
}

type UpdateMovieInput struct {
	Title       string `json:"title"`
	Director    string `json:"director"`
	Year        int    `json:"year"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
}

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	config.DB.Preload("Reviews").Find(&movies)
	c.JSON(http.StatusOK, movies)
}

func GetMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.Preload("Reviews").First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func CreateMovie(c *gin.Context) {
	var input CreateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := models.Movie{
		Title:       input.Title,
		Director:    input.Director,
		Year:        input.Year,
		Genre:       input.Genre,
		Description: input.Description,
	}

	config.DB.Create(&movie)
	c.JSON(http.StatusCreated, movie)
}

func UpdateMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	var input UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&movie).Updates(input)
	c.JSON(http.StatusOK, movie)
}

func DeleteMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	config.DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
