package service

import (
	"movies-api/models"
	"movies-api/repository"
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
