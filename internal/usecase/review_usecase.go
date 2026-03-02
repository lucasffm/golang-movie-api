package usecase

import (
	"movie-api/internal/domain"
	"time"
)

type reviewUseCase struct {
	repo      domain.ReviewRepository
	movieRepo domain.MovieRepository
}

func NewReviewUseCase(repo domain.ReviewRepository, movieRepo domain.MovieRepository) domain.ReviewUseCase {
	return &reviewUseCase{repo: repo, movieRepo: movieRepo}
}

func (uc *reviewUseCase) GetReviews() ([]domain.Review, error) {
	return uc.repo.FindAll()
}

func (uc *reviewUseCase) GetReview(id uint) (*domain.Review, error) {
	return uc.repo.FindByID(id)
}

func (uc *reviewUseCase) GetMovieReviews(movieID uint) ([]domain.Review, error) {
	return uc.repo.FindByMovieID(movieID)
}

func (uc *reviewUseCase) CreateReview(input domain.CreateReviewInput) (*domain.Review, error) {
	if _, err := uc.movieRepo.FindByID(input.MovieID); err != nil {
		return nil, domain.ErrMovieNotFound
	}

	watchedAt := input.WatchedAt
	if watchedAt.IsZero() {
		watchedAt = time.Now()
	}

	review := &domain.Review{
		MovieID:   input.MovieID,
		Rating:    input.Rating,
		Comment:   input.Comment,
		WatchedAt: watchedAt,
	}

	if err := uc.repo.Create(review); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(review.ID)
}

func (uc *reviewUseCase) UpdateReview(id uint, input domain.UpdateReviewInput) (*domain.Review, error) {
	review, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if input.Rating != 0 {
		data["rating"] = input.Rating
	}
	if input.Comment != "" {
		data["comment"] = input.Comment
	}
	if !input.WatchedAt.IsZero() {
		data["watched_at"] = input.WatchedAt
	}

	if err := uc.repo.Update(review, data); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(id)
}

func (uc *reviewUseCase) DeleteReview(id uint) error {
	review, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(review)
}
