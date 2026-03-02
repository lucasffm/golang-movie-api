package usecase_test

import (
	"errors"
	"movie-api/internal/domain"
	"movie-api/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) FindAll() ([]domain.Movie, error) {
	args := m.Called()
	return args.Get(0).([]domain.Movie), args.Error(1)
}

func (m *MockMovieRepository) FindByID(id uint) (*domain.Movie, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Movie), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieRepository) Create(movie *domain.Movie) error {
	args := m.Called(movie)
	if args.Error(0) == nil {
		movie.ID = 1 // simulate DB assigning ID
	}
	return args.Error(0)
}

func (m *MockMovieRepository) Update(movie *domain.Movie, data map[string]interface{}) error {
	args := m.Called(movie, data)
	return args.Error(0)
}

func (m *MockMovieRepository) Delete(movie *domain.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func setup() (*MockMovieRepository, domain.MovieUseCase) {
	repo := new(MockMovieRepository)
	uc := usecase.NewMovieUseCase(repo)
	return repo, uc
}

func TestGetMovies(t *testing.T) {
	repo, uc := setup()

	expectedMovies := []domain.Movie{
		{ID: 1, Title: "Movie 1"},
		{ID: 2, Title: "Movie 2"},
	}

	repo.On("FindAll").Return(expectedMovies, nil)

	movies, err := uc.GetMovies()

	assert.NoError(t, err)
	assert.Equal(t, expectedMovies, movies)
	repo.AssertExpectations(t)
}

func TestGetMovie_Success(t *testing.T) {
	repo, uc := setup()

	expectedMovie := &domain.Movie{ID: 1, Title: "Movie 1"}
	repo.On("FindByID", uint(1)).Return(expectedMovie, nil)

	movie, err := uc.GetMovie(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMovie, movie)
	repo.AssertExpectations(t)
}

func TestGetMovie_Error(t *testing.T) {
	repo, uc := setup()

	repo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	movie, err := uc.GetMovie(1)

	assert.Error(t, err)
	assert.Nil(t, movie)
	repo.AssertExpectations(t)
}

func TestCreateMovie_Success(t *testing.T) {
	repo, uc := setup()

	input := domain.CreateMovieInput{
		Title:       "New Movie",
		Director:    "Director",
		Year:        2024,
		Genre:       "Action",
		Description: "A great movie",
	}

	repo.On("Create", mock.Anything).Return(nil)

	movie, err := uc.CreateMovie(input)

	assert.NoError(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, uint(1), movie.ID)
	assert.Equal(t, input.Title, movie.Title)
	repo.AssertExpectations(t)
}

func TestCreateMovie_Error(t *testing.T) {
	repo, uc := setup()

	input := domain.CreateMovieInput{
		Title: "New Movie",
	}

	repo.On("Create", mock.Anything).Return(errors.New("db error"))

	movie, err := uc.CreateMovie(input)

	assert.Error(t, err)
	assert.Nil(t, movie)
	repo.AssertExpectations(t)
}

func TestUpdateMovie_Success(t *testing.T) {
	repo, uc := setup()

	existingMovie := &domain.Movie{ID: 1, Title: "Old Title", Director: "Old Director", Year: 2000, Genre: "Old Genre", Description: "Old Description"}
	updatedMovie := &domain.Movie{ID: 1, Title: "New Title", Director: "New Director", Year: 2001, Genre: "New Genre", Description: "New Description"}
	input := domain.UpdateMovieInput{Title: "New Title", Director: "New Director", Year: 2001, Genre: "New Genre", Description: "New Description"}

	repo.On("FindByID", uint(1)).Return(existingMovie, nil).Once()
	repo.On("Update", existingMovie, mock.Anything).Return(nil)
	repo.On("FindByID", uint(1)).Return(updatedMovie, nil).Once()

	movie, err := uc.UpdateMovie(1, input)

	assert.NoError(t, err)
	assert.Equal(t, updatedMovie, movie)
	repo.AssertExpectations(t)
}

func TestUpdateMovie_NotFound(t *testing.T) {
	repo, uc := setup()

	input := domain.UpdateMovieInput{Title: "New Title"}
	repo.On("FindByID", uint(1)).Return(nil, domain.ErrMovieNotFound)

	movie, err := uc.UpdateMovie(1, input)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMovieNotFound, err)
	assert.Nil(t, movie)
	repo.AssertExpectations(t)
}

func TestUpdateMovie_UpdateError(t *testing.T) {
	repo, uc := setup()

	existingMovie := &domain.Movie{ID: 1, Title: "Old Title"}
	input := domain.UpdateMovieInput{Title: "New Title"}

	repo.On("FindByID", uint(1)).Return(existingMovie, nil).Once()
	repo.On("Update", existingMovie, mock.Anything).Return(errors.New("update error"))

	movie, err := uc.UpdateMovie(1, input)

	assert.Error(t, err)
	assert.Nil(t, movie)
	repo.AssertExpectations(t)
}

func TestDeleteMovie_Success(t *testing.T) {
	repo, uc := setup()

	existingMovie := &domain.Movie{ID: 1, Title: "Movie 1"}

	repo.On("FindByID", uint(1)).Return(existingMovie, nil)
	repo.On("Delete", existingMovie).Return(nil)

	err := uc.DeleteMovie(1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestDeleteMovie_NotFound(t *testing.T) {
	repo, uc := setup()

	repo.On("FindByID", uint(1)).Return(nil, domain.ErrMovieNotFound)

	err := uc.DeleteMovie(1)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrMovieNotFound, err)
	repo.AssertExpectations(t)
}
