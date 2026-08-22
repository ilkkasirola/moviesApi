package service

import (
	// "errors"
	"database/sql"
	stdErr "errors"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/repository"
	"strings"
)

// var ErrMovieExists = errors.New("movie already exists")

type MovieService struct {
	MRepo *repository.MovieRepository
	ARepo *repository.ActorRepository
	GRepo *repository.GenreRepository
}

func NewMovieService(mRepo *repository.MovieRepository, aRepo *repository.ActorRepository, gRepo *repository.GenreRepository) *MovieService {
	return &MovieService{
		MRepo: mRepo,
		ARepo: aRepo,
		GRepo: gRepo,
	}
}

func (s *MovieService) CreateMovie(m models.Movie) (int64, error) {
	movieID, err := s.MRepo.Create(
		m.Title,
		m.ReleaseYear,
		m.Duration,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, appErr.ErrMovieExists
		}
		return 0, err
	}

	for _, actor := range m.Actors {
		_, err = s.ARepo.CreateActor(actor.Name, actor.BirthDate)
		if err != nil {
			return 0, err
		}
		actorID, err := s.ARepo.GetActorID(actor.Name, actor.BirthDate)
		if err != nil {
			return 0, err
		}
		err = s.ARepo.LinkActor(movieID, actorID)
		if err != nil {
			return 0, err
		}

		err = s.ARepo.LinkActor(movieID, actorID)
		if err != nil {
			return 0, err
		}
	}

	return movieID, nil
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	return s.MRepo.GetAllMovies()
}

func (s *MovieService) GetMovieByID(id int) (models.Movie, error) {
	movie, err := s.MRepo.GetMovieByID(id)
	if stdErr.Is(err, sql.ErrNoRows) {
		return models.Movie{}, appErr.ErrMovieNotFound
	}
	if err != nil {
		return models.Movie{}, err
	}
	return movie, nil
}
