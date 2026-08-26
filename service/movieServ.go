package service

import (
	// "errors"
	"context"
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

func (s *MovieService) SearchMovieByGenreID(genreID int) ([]models.Movie, error) {
	movies, err := s.MRepo.SearchMovieByGenreID(genreID)
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

func (s *MovieService) SearchMovieByReleaseYear(releaseYear int) ([]models.Movie, error) {
	movies, err := s.MRepo.SearchMovieByReleaseYear(releaseYear)
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
	genres, err := s.GRepo.GetGenreByMovieID(movie.ID)
	if err != nil {
		return models.Movie{}, err
	}
	movie.Genres = genres
	actors, err := s.ARepo.GetByMovieID(movie.ID)
	if err != nil {
		return models.Movie{}, err
	}
	movie.Actors = actors
	return movie, nil
}

func (s *MovieService) GetMoviesByActorID(actorID int) ([]models.Movie, error) {
	return s.ARepo.GetMoviesByActorID(actorID)
}

func (s *MovieService) DeleteMovie(id int, force bool) error {
	return s.MRepo.DeleteMovie(id, force)
}

func (s *MovieService) PatchMovie(ctx context.Context, id int, title *string, year, duration *int, genreIds, actorIds *[]int) (*models.Movie, error) {
	if genreIds != nil {
		for _, g := range *genreIds {
			if _, err := s.GRepo.GetByID(g); err != nil {
				if stdErr.Is(err, sql.ErrNoRows) {
					return nil, appErr.ErrGenreNotFound
				}
				return nil, err
			}
		}
	}
	if actorIds != nil {
		for _, a := range *actorIds {
			if _, err := s.ARepo.GetByID(a); err != nil {
				if stdErr.Is(err, sql.ErrNoRows) {
					return nil, appErr.ErrActorNotFound
				}
				return nil, err
			}
		}
	}

	movie, err := s.MRepo.PatchMovie(ctx, id, title, year, duration, genreIds, actorIds)
	if err != nil {
		return nil, err
	}
	genres, err := s.GRepo.GetGenreByMovieID(movie.ID)
	if err != nil {
		return nil, err
	}
	movie.Genres = genres
	actors, err := s.ARepo.GetByMovieID(movie.ID)
	if err != nil {
		return nil, err
	}
	movie.Actors = actors
	return movie, nil
}
