package repository

import (
	"errors"
	"movie-api/internal/domain"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) FindAll() ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Preload("Movie").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) FindByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("Movie").First(&review, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrReviewNotFound
		}
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByMovieID(movieID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Where("movie_id = ?", movieID).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) Update(review *domain.Review, data map[string]interface{}) error {
	return r.db.Model(review).Updates(data).Error
}

func (r *reviewRepository) Delete(review *domain.Review) error {
	return r.db.Delete(review).Error
}
