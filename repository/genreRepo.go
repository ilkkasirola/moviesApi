package repository

import (
	"context"
	"database/sql"
	appErr "movies-api/errors"
	"movies-api/models"
)

type GenreRepository struct {
	DB *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{DB: db}
}

func (r *GenreRepository) GetAll() ([]models.Genre, error) {
	rows, err := r.DB.Query(`
		SELECT id, name From Genre
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, rows.Err()
}

func (r *GenreRepository) GetByID(id int) (models.Genre, error) {
	var g models.Genre

	err := r.DB.QueryRow(`
		SELECT id, name FROM Genre WHERE id = ?
		`, id).Scan(&g.ID, &g.Name)

	return g, err
}

func (r *GenreRepository) CreateGenre(name string) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT OR IGNORE INTO Genre (name) VALUES (?)", name,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *GenreRepository) GetGenreID(name string) (int, error) {
	var id int
	err := r.DB.QueryRow("SELECT id FROM Genre WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *GenreRepository) DeleteGenre(id int, force bool) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var name string
	if err := tx.QueryRow(`SELECT name FROM Genre WHERE id = ?`, id).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return appErr.ErrGenreNotFound
		}
		return err
	}
	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM MovieGenre WHERE genreId = ?`, id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		if !force {
			return appErr.NewGenreInUse(name, count)
		}
		if _, err := tx.Exec(`DELETE FROM MovieGenre WHERE genreId = ?`, id); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM Genre WHERE id = ?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *GenreRepository) PatchGenre(ctx context.Context, id int, name string) (*models.Genre, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, "UPDATE Genre SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return nil, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, appErr.ErrGenreNotFound
	}

	var g models.Genre
	row := tx.QueryRowContext(ctx, "SELECT id, name FROM Genre WHERE id = ?", id)
	if err := row.Scan(&g.ID, &g.Name); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GenreRepository) GetGenreByMovieID(movieID int) ([]models.Genre, error) {
	rows, err := r.DB.Query(`
		SELECT g.id, g.name FROM Genre g
		INNER JOIN MovieGenre mg ON g.id = mg.GenreId WHERE mg.movieId = ?
		`, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, rows.Err()
}
