package repository

import (
	"errors"
	"movie-api/internal/domain"

	"gorm.io/gorm"
)

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) domain.MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) FindAll() ([]domain.Movie, error) {
	var movies []domain.Movie
	err := r.db.Preload("Reviews").Find(&movies).Error
	return movies, err
}

func (r *movieRepository) FindByID(id uint) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.db.Preload("Reviews").First(&movie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrMovieNotFound
		}
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Create(movie *domain.Movie) error {
	return r.db.Create(movie).Error
}

func (r *movieRepository) Update(movie *domain.Movie, data map[string]interface{}) error {
	return r.db.Model(movie).Updates(data).Error
}

func (r *movieRepository) Delete(movie *domain.Movie) error {
	return r.db.Delete(movie).Error
}
