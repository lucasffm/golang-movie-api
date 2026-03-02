package handler

import (
	"errors"
	"movie-api/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	useCase domain.MovieUseCase
}

func NewMovieHandler(useCase domain.MovieUseCase) *MovieHandler {
	return &MovieHandler{useCase: useCase}
}

func (h *MovieHandler) GetMovies(c *gin.Context) {
	movies, err := h.useCase.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	movie, err := h.useCase.GetMovie(id)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var input domain.CreateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie, err := h.useCase.CreateMovie(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input domain.UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie, err := h.useCase.UpdateMovie(id, input)
	if err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.useCase.DeleteMovie(id); err != nil {
		if errors.Is(err, domain.ErrMovieNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
