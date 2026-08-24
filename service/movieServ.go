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
	movieID, err := s.MRepo.CreateMovie(
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
	}

	for _, genre := range m.Genres {
		_, err = s.GRepo.CreateGenre(genre.Name)
		if err != nil {
			return 0, err
		}
		genreID, err := s.GRepo.GetGenreID(genre.Name)
		if err != nil {
			return 0, err
		}
		err = s.GRepo.LinkGenre(movieID, genreID)
		if err != nil {
			return 0, err
		}
	}

	return movieID, nil
}

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	movies, err := s.MRepo.GetAllMovies()
	if err != nil {
		return nil, err
	}

	var res []models.Movie

	for _, m := range movies {
		actors, err := s.ARepo.GetByMovieID(m.ID)
		if err != nil {
			return nil, err
		}
		genres, err := s.GRepo.GetGenreByMovieID(m.ID)
		if err != nil {
			return nil, err
		}
		res = append(res, models.Movie{
			ID:          m.ID,
			Title:       m.Title,
			ReleaseYear: m.ReleaseYear,
			Duration:    m.Duration,
			Actors:      actors,
			Genres:      genres,
		})
	}
	return res, nil
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
