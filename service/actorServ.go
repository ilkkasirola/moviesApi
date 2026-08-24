package service

import (
	"context"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/repository"
	"strings"
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

func (s *ActorService) DeleteActor(id int) error {
	return s.Repo.DeleteActor(id)
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
func (s *ActorService) GetActorsByMovieID(movieID int) ([]models.Actor, error) {
	actor, err := s.Repo.GetByMovieID(movieID)
	if err != nil {
		return nil, err
	}
	return actor, nil
}

func (s *ActorService) CreateActor(a models.Actor) (int64, error) {
	actorID, err := s.Repo.CreateActor(a.Name, a.BirthDate)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, appErr.ErrActorExists
		}
		return 0, err
	}
	return actorID, nil
}

func (s *ActorService) PatchActor(ctx context.Context, id int, name string, birthDate string) (*models.Actor, error) {
	name = strings.TrimSpace(name)
	return s.Repo.PatchActor(ctx, id, name, birthDate)
}
