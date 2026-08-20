package service

import (
	// "errors"

	"movies-api/errors"
	"movies-api/models"
	"movies-api/repository"
	"strings"
)

// var ErrMovieExists = errors.New("movie already exists")

type MovieService struct {
	Repo *repository.MovieRepository
}

func NewMovieService(r *repository.MovieRepository) *MovieService {
	return &MovieService{Repo: r}
}

func (s *MovieService) CreateMovie(m models.Movie) (int64, error) {
	movieID, err := s.Repo.Create(
		m.Title,
		m.ReleaseYear,
		m.Duration,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, errors.ErrMovieExists
		}
		return 0, err
	}

	for _, actor := range m.Actors {
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

		// 	actorID, err := s.Repo.GetActorID(actor.Name, actor.BirthDate.Time)
		// 	if err != nil {
		// 		return 0, err
		// 	}

		// 	err = s.Repo.LinkActor(movieID, actorID)
		// 	if err != nil {
		// 		return 0, err
		// 	}
		// }
	}
	return movieID, nil
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	movies, err := s.Repo.GetAllMovies()
	if err != nil {
		return nil, err
	}

	for i := range movies {
		actors, _ := s.Repo.GetActor
	}
}
