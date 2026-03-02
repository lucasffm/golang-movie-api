package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrMovieNotFound = errors.New("movie not found")

type Movie struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Director    string         `json:"director"`
	Year        int            `json:"year"`
	Genre       string         `json:"genre"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Reviews     []Review       `json:"reviews,omitempty"`
}

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

type MovieRepository interface {
	FindAll() ([]Movie, error)
	FindByID(id uint) (*Movie, error)
	Create(movie *Movie) error
	Update(movie *Movie, data map[string]interface{}) error
	Delete(movie *Movie) error
}

type MovieUseCase interface {
	GetMovies() ([]Movie, error)
	GetMovie(id uint) (*Movie, error)
	CreateMovie(input CreateMovieInput) (*Movie, error)
	UpdateMovie(id uint, input UpdateMovieInput) (*Movie, error)
	DeleteMovie(id uint) error
}
