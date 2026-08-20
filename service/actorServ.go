package service

import (
	"movies-api/models"
	"movies-api/repository"
)

type ActorService struct {
	Repo *repository.ActorRepository
}

func NewActorService(r *repository.ActorRepository) *ActorService {
	return &ActorService{Repo: r}
}

func (s *ActorService) GetAllActors() ([]models.Actor, error) {
	return s.Repo.GetAll()
}

func (s *ActorService) GetActorByID(id int) (models.Actor, []models.Movie, error) {
	actor, err := s.Repo.GetByID(id)
	if err != nil {
		return actor, nil, err
	}

	movies, err := s.Repo.GetMoviesByActorID(id)
	if err != nil {
		return actor, nil, err
	}
	return actor, movies, nil
}
