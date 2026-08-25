package repository

import (
	"context"
	"database/sql"
	aErr "movies-api/errors"
	"movies-api/models"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (r *MovieRepository) GetAllMovies() ([]models.Movie, error) {
	rows, err := r.DB.Query(`
		SELECT id, title, releaseYear, duration FROM Movie
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var m models.Movie
		err := rows.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, rows.Err()
}

func (r *MovieRepository) GetMovieByID(id int) (models.Movie, error) {
	var m models.Movie

	err := r.DB.QueryRow(`
		SELECT id, title, releaseYear, duration FROM Movie WHERE id = ?`,
		id).Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration)

	return m, err
}

func (r *MovieRepository) CreateMovie(title string, year int, duration int) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT INTO Movie (title, releaseYear, duration) VALUES (?, ?, ?)",
		title, year, duration,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MovieRepository) DeleteMovie(id int) error {
	_, err := r.DB.Exec(`DELETE FROM movie WHERE id = ?`, id)
	return err
}

func (r *MovieRepository) PatchMovie(ctx context.Context, id int, title *string, year, duration *int, genreIds, actorIds *[]int) (*models.Movie, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE Movie
		SET title = COALESCE(?, title),
			releaseYear = COALESCE(?, releaseYear),
			duration = COALESCE(?, duration)
		WHERE id = ?`,
		title, year, duration, id,
	)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, aErr.ErrMovieNotFound
	}
	if genreIds != nil {
		_, err = tx.ExecContext(
			ctx, `
			DELETE FROM MovieGenre
			WHERE movieId = ?
			`, id,
		)
		if err != nil {
			return nil, err
		}

		for _, v := range *genreIds {
			_, err := tx.ExecContext(ctx, `
				INSERT OR IGNORE INTO MovieGenre (movieID, genreID) VALUES (?, ?)`,
				id, v,
			)
			if err != nil {
				return nil, err
			}
		}

	}
	if actorIds != nil {
		_, err = tx.ExecContext(
			ctx, `
			DELETE FROM MovieActor
			WHERE movieId = ?
			`, id,
		)
		if err != nil {
			return nil, err
		}

		for _, v := range *actorIds {
			_, err := tx.ExecContext(ctx, `
				INSERT OR IGNORE INTO MovieActor (movieID, actorID) VALUES (?, ?)`,
				id, v,
			)
			if err != nil {
				return nil, err
			}
		}
	}
	var m models.Movie
	row := tx.QueryRowContext(ctx, `SELECT id, title, releaseYear, duration FROM Movie WHERE id = ?`, id)
	if err := row.Scan(&m.ID, &m.Title, &m.ReleaseYear, &m.Duration); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &m, nil
}
