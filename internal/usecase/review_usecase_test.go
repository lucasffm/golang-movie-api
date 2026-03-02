package usecase_test

import (
	"errors"
	"movie-api/internal/domain"
	"movie-api/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) FindAll() ([]domain.Review, error) {
	args := m.Called()
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepository) FindByID(id uint) (*domain.Review, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Review), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockReviewRepository) FindByMovieID(movieID uint) ([]domain.Review, error) {
	args := m.Called(movieID)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepository) Create(review *domain.Review) error {
	args := m.Called(review)
	if args.Error(0) == nil {
		review.ID = 1
	}
	return args.Error(0)
}

func (m *MockReviewRepository) Update(review *domain.Review, data map[string]interface{}) error {
	args := m.Called(review, data)
	return args.Error(0)
}

func (m *MockReviewRepository) Delete(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func setupReview() (*MockReviewRepository, *MockMovieRepository, domain.ReviewUseCase) {
	repo := new(MockReviewRepository)
	movieRepo := new(MockMovieRepository)
	uc := usecase.NewReviewUseCase(repo, movieRepo)
	return repo, movieRepo, uc
}

func TestGetReviews(t *testing.T) {
	repo, _, uc := setupReview()

	expectedReviews := []domain.Review{
		{ID: 1, MovieID: 1, Rating: 8, Comment: "Great movie"},
		{ID: 2, MovieID: 2, Rating: 7, Comment: "Good movie"},
	}

	repo.On("FindAll").Return(expectedReviews, nil)

	reviews, err := uc.GetReviews()

	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)
	repo.AssertExpectations(t)
}

func TestGetReview_Success(t *testing.T) {
	repo, _, uc := setupReview()

	expectedReview := &domain.Review{ID: 1, MovieID: 1, Rating: 8, Comment: "Great movie"}
	repo.On("FindByID", uint(1)).Return(expectedReview, nil)

	review, err := uc.GetReview(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedReview, review)
	repo.AssertExpectations(t)
}

func TestGetReview_Error(t *testing.T) {
	repo, _, uc := setupReview()

	repo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	review, err := uc.GetReview(1)

	assert.Error(t, err)
	assert.Nil(t, review)
	repo.AssertExpectations(t)
}

func TestGetMovieReviews(t *testing.T) {
	repo, _, uc := setupReview()

	expectedReviews := []domain.Review{
		{ID: 1, MovieID: 1, Rating: 8, Comment: "Great movie"},
		{ID: 2, MovieID: 1, Rating: 7, Comment: "Good movie"},
	}

	repo.On("FindByMovieID", uint(1)).Return(expectedReviews, nil)

	reviews, err := uc.GetMovieReviews(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)
	repo.AssertExpectations(t)
}

func TestCreateReview_Success(t *testing.T) {
	repo, movieRepo, uc := setupReview()

	movie := &domain.Movie{ID: 1, Title: "Movie 1"}
	input := domain.CreateReviewInput{
		MovieID: 1,
		Rating:  8,
		Comment: "Great movie",
	}

	movieRepo.On("FindByID", uint(1)).Return(movie, nil)
	repo.On("Create", mock.Anything).Return(nil)
	repo.On("FindByID", uint(1)).Return(&domain.Review{ID: 1, MovieID: 1, Rating: 8, Comment: "Great movie"}, nil)

	review, err := uc.CreateReview(input)

	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, uint(1), review.ID)
	assert.Equal(t, input.Rating, review.Rating)
	movieRepo.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateReview_MovieNotFound(t *testing.T) {
	_, movieRepo, uc := setupReview()

	input := domain.CreateReviewInput{
		MovieID: 1,
		Rating:  8,
		Comment: "Great movie",
	}

	movieRepo.On("FindByID", uint(1)).Return(nil, domain.ErrMovieNotFound)

	review, err := uc.CreateReview(input)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMovieNotFound, err)
	assert.Nil(t, review)
	movieRepo.AssertExpectations(t)
}

func TestCreateReview_DefaultWatchedAt(t *testing.T) {
	repo, movieRepo, uc := setupReview()

	movie := &domain.Movie{ID: 1, Title: "Movie 1"}
	input := domain.CreateReviewInput{
		MovieID: 1,
		Rating:  8,
		Comment: "Great movie",
	}

	movieRepo.On("FindByID", uint(1)).Return(movie, nil)
	repo.On("Create", mock.MatchedBy(func(r *domain.Review) bool {
		return !r.WatchedAt.IsZero()
	})).Return(nil)
	repo.On("FindByID", uint(1)).Return(&domain.Review{ID: 1, MovieID: 1, Rating: 8, Comment: "Great movie", WatchedAt: time.Now()}, nil)

	review, err := uc.CreateReview(input)

	assert.NoError(t, err)
	assert.NotNil(t, review)
	movieRepo.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCreateReview_Error(t *testing.T) {
	repo, movieRepo, uc := setupReview()

	movie := &domain.Movie{ID: 1, Title: "Movie 1"}
	input := domain.CreateReviewInput{
		MovieID: 1,
		Rating:  8,
		Comment: "Great movie",
	}

	movieRepo.On("FindByID", uint(1)).Return(movie, nil)
	repo.On("Create", mock.Anything).Return(errors.New("db error"))

	review, err := uc.CreateReview(input)

	assert.Error(t, err)
	assert.Nil(t, review)
	movieRepo.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestUpdateReview_Success(t *testing.T) {
	repo, _, uc := setupReview()

	existingReview := &domain.Review{ID: 1, MovieID: 1, Rating: 5, Comment: "Old comment"}
	updatedReview := &domain.Review{ID: 1, MovieID: 1, Rating: 8, Comment: "New comment"}
	input := domain.UpdateReviewInput{Rating: 8, Comment: "New comment"}

	repo.On("FindByID", uint(1)).Return(existingReview, nil).Once()
	repo.On("Update", existingReview, mock.Anything).Return(nil)
	repo.On("FindByID", uint(1)).Return(updatedReview, nil).Once()

	review, err := uc.UpdateReview(1, input)

	assert.NoError(t, err)
	assert.Equal(t, updatedReview, review)
	repo.AssertExpectations(t)
}

func TestUpdateReview_NotFound(t *testing.T) {
	repo, _, uc := setupReview()

	input := domain.UpdateReviewInput{Rating: 8}
	repo.On("FindByID", uint(1)).Return(nil, domain.ErrReviewNotFound)

	review, err := uc.UpdateReview(1, input)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrReviewNotFound, err)
	assert.Nil(t, review)
	repo.AssertExpectations(t)
}

func TestUpdateReview_UpdateError(t *testing.T) {
	repo, _, uc := setupReview()

	existingReview := &domain.Review{ID: 1, MovieID: 1, Rating: 5}
	input := domain.UpdateReviewInput{Rating: 8}

	repo.On("FindByID", uint(1)).Return(existingReview, nil).Once()
	repo.On("Update", existingReview, mock.Anything).Return(errors.New("update error"))

	review, err := uc.UpdateReview(1, input)

	assert.Error(t, err)
	assert.Nil(t, review)
	repo.AssertExpectations(t)
}

func TestDeleteReview_Success(t *testing.T) {
	repo, _, uc := setupReview()

	existingReview := &domain.Review{ID: 1, MovieID: 1, Rating: 8}

	repo.On("FindByID", uint(1)).Return(existingReview, nil)
	repo.On("Delete", existingReview).Return(nil)

	err := uc.DeleteReview(1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteReview_NotFound(t *testing.T) {
	repo, _, uc := setupReview()

	repo.On("FindByID", uint(1)).Return(nil, domain.ErrReviewNotFound)

	err := uc.DeleteReview(1)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrReviewNotFound, err)
	repo.AssertExpectations(t)
}
