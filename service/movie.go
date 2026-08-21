package service

import (
	"errors"
	"movies-api/models"
	"movies-api/repository"
	"strings"
)

var ErrMovieExists = errors.New("movie already exists")

type MovieService struct {
	Repo *repository.MovieRepository
}

func NewMovieService(repo *repository.MovieRepository) *MovieService {
	return &MovieService{Repo: repo}
}

func (s *MovieService) CreateMovie(movie models.Movie) (int64, error) {
	movieID, err := s.Repo.Create(
		movie.Title,
		movie.ReleaseYear,
		movie.Duration,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, ErrMovieExists
		}
		return 0, err
	}

	for _, actor := range movie.Actors {
		err = s.Repo.CreateActor(actor.Name, actor.BirthDate.Time)
		if err != nil {
			return 0, err
		}

		actorID, err := s.Repo.GetActorID(actor.Name, actor.BirthDate.Time)
		if err != nil {
			return 0, err
		}

		err = s.Repo.LinkActor(movieID, actorID)
		if err != nil {
			return 0, err
		}
	}

	return movieID, nil
}
