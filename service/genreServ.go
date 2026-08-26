package service

import (
	"context"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/repository"
	"strings"
)

type GenreService struct {
	Repo *repository.GenreRepository
}

func NewGenreService(r *repository.GenreRepository) *GenreService {
	return &GenreService{Repo: r}
}

func (s *GenreService) GetAllGenres() ([]models.Genre, error) {
	return s.Repo.GetAll()
}

func (s *GenreService) GetGenreByID(id int) (models.Genre, error) {
	genre, err := s.Repo.GetByID(id)
	if err != nil {
		return models.Genre{}, err
	}
	return genre, nil

}

func (s *GenreService) CreateGenre(g models.Genre) (int64, error) {
	genreID, err := s.Repo.CreateGenre(
		g.Name,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, appErr.ErrGenreExists
		}
		return 0, err
	}

	return genreID, nil
}

func (s *GenreService) DeleteGenre(id int, force bool) error {
	return s.Repo.DeleteGenre(id, force)
}

func (s *GenreService) PatchGenre(ctx context.Context, id int, name string) (*models.Genre, error) {
	name = strings.TrimSpace(name)
	return s.Repo.PatchGenre(ctx, id, name)
}

func (s *GenreService) GetGenreByMovieID(id int) ([]models.Genre, error) {
	return s.Repo.GetGenreByMovieID(id)
}
