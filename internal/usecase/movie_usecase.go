package usecase

import "movie-api/internal/domain"

type movieUseCase struct {
	repo domain.MovieRepository
}

func NewMovieUseCase(repo domain.MovieRepository) domain.MovieUseCase {
	return &movieUseCase{repo: repo}
}

func (uc *movieUseCase) GetMovies() ([]domain.Movie, error) {
	return uc.repo.FindAll()
}

func (uc *movieUseCase) GetMovie(id uint) (*domain.Movie, error) {
	return uc.repo.FindByID(id)
}

func (uc *movieUseCase) CreateMovie(input domain.CreateMovieInput) (*domain.Movie, error) {
	movie := &domain.Movie{
		Title:       input.Title,
		Director:    input.Director,
		Year:        input.Year,
		Genre:       input.Genre,
		Description: input.Description,
	}
	if err := uc.repo.Create(movie); err != nil {
		return nil, err
	}
	return movie, nil
}

func (uc *movieUseCase) UpdateMovie(id uint, input domain.UpdateMovieInput) (*domain.Movie, error) {
	movie, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	if input.Title != "" {
		data["title"] = input.Title
	}
	if input.Director != "" {
		data["director"] = input.Director
	}
	if input.Year != 0 {
		data["year"] = input.Year
	}
	if input.Genre != "" {
		data["genre"] = input.Genre
	}
	if input.Description != "" {
		data["description"] = input.Description
	}

	if err := uc.repo.Update(movie, data); err != nil {
		return nil, err
	}
	return uc.repo.FindByID(id)
}

func (uc *movieUseCase) DeleteMovie(id uint) error {
	movie, err := uc.repo.FindByID(id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(movie)
}
