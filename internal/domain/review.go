package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrReviewNotFound = errors.New("review not found")

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MovieID   uint           `json:"movie_id" gorm:"not null;index"`
	Rating    int            `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 10"`
	Comment   string         `json:"comment"`
	WatchedAt time.Time      `json:"watched_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Movie     Movie          `json:"movie,omitempty"`
}

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

type ReviewRepository interface {
	FindAll() ([]Review, error)
	FindByID(id uint) (*Review, error)
	FindByMovieID(movieID uint) ([]Review, error)
	Create(review *Review) error
	Update(review *Review, data map[string]interface{}) error
	Delete(review *Review) error
}

type ReviewUseCase interface {
	GetReviews() ([]Review, error)
	GetReview(id uint) (*Review, error)
	GetMovieReviews(movieID uint) ([]Review, error)
	CreateReview(input CreateReviewInput) (*Review, error)
	UpdateReview(id uint, input UpdateReviewInput) (*Review, error)
	DeleteReview(id uint) error
}
