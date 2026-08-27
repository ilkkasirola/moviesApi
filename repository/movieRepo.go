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

func (r *MovieRepository) SearchMovieByGenreID(genreID int) ([]models.Movie, error) {
	rows, err := r.DB.Query(`
		SELECT m.id, m.title, m.releaseYear, m.duration FROM Movie m
		INNER JOIN MovieGenre mg ON m.id = mg.movieId
		WHERE mg.genreId = ?`, genreID)
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
func (r *MovieRepository) SearchMovieByReleaseYear(releaseYear int) ([]models.Movie, error) {
	rows, err := r.DB.Query(`
		SELECT m.id, m.title, m.releaseYear, m.duration FROM Movie m
		WHERE m.releaseYear = ?`, releaseYear)
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

func (r *MovieRepository) CreateMovie(m models.Movie) (int64, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	res, err := tx.Exec(
		"INSERT INTO Movie (title, releaseYear, duration) VALUES (?, ?, ?)",
		m.Title, m.ReleaseYear, m.Duration,
	)
	if err != nil {
		return 0, err
	}
	movieID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, actor := range m.Actors {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO Actor (name, birthDate) VALUES (?, ?)`,
			actor.Name, actor.BirthDate,
		); err != nil {
			return 0, err
		}
		var actorID int64
		if err := tx.QueryRow(
			`SELECT id FROM Actor WHERE name = ? AND birthDate = ?`,
			actor.Name, actor.BirthDate,
		).Scan(&actorID); err != nil {
			return 0, err
		}
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO MovieActor (movieId, actorId) VALUES (?, ?)`,
			movieID, actorID,
		); err != nil {
			return 0, err
		}
	}

	for _, genre := range m.Genres {
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO Genre (name) VALUES (?)`,
			genre.Name,
		); err != nil {
			return 0, err
		}
		var genreID int64
		if err := tx.QueryRow(
			`SELECT id FROM Genre WHERE name = ?`,
			genre.Name,
		).Scan(&genreID); err != nil {
			return 0, err
		}
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO MovieGenre (movieId, genreId) VALUES (?, ?)`,
			movieID, genreID,
		); err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return movieID, nil
}

func (r *MovieRepository) DeleteMovie(id int, force bool) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var title string
	if err := tx.QueryRow(`SELECT title FROM Movie WHERE id = ?`, id).Scan(&title); err != nil {
		if err == sql.ErrNoRows {
			return aErr.ErrMovieNotFound
		}
		return err
	}

	var actorCount, genreCount int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM MovieActor WHERE movieId = ?`, id).Scan(&actorCount); err != nil {
		return err
	}
	if err := tx.QueryRow(`SELECT COUNT(*) FROM MovieGenre WHERE movieId = ?`, id).Scan(&genreCount); err != nil {
		return err
	}
	if actorCount > 0 || genreCount > 0 {
		if !force {
			return aErr.NewMovieInUse(title, actorCount, genreCount)
		}
		if _, err := tx.Exec(`DELETE FROM MovieActor WHERE movieId = ?`, id); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM MovieGenre WHERE movieId = ?`, id); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM Movie WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *MovieRepository) CountLinkedActors(id int) (int, error) {
	var count int
	err := r.DB.QueryRow(
		"SELECT COUNT(*) FROM MovieActor WHERE movieId = ?", id,
	).Scan(&count)
	return count, err
}

// func (r *MovieRepository) CountLinkedGenres(id int) (int, error) {
// 	var count int
// 	err := r.DB.QueryRow(
// 		"SELECT COUNT(*) FROM MovieGenre WHERE movieId = ?", id,
// 	).Scan(&count)
// 	return count, err
// }

// func (r *MovieRepository) UnlinkAllActorsAndGenres(id int) error {
// 	if _, err := r.DB.Exec("DELETE FROM MovieActor WHERE movieId = ?", id); err != nil {
// 		return err
// 	}
// 	_, err := r.DB.Exec("DELETE FROM MovieGenre WHERE movieId = ?", id)
// 	return err
// }

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
				INSERT OR IGNORE INTO MovieGenre (movieId, genreId) VALUES (?, ?)`,
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
				INSERT OR IGNORE INTO MovieActor (movieId, actorId) VALUES (?, ?)`,
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
