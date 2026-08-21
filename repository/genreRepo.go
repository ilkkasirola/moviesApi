package repository

import (
	"database/sql"
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
		"INSERT INTO Genre (name) VALUES (?)", name,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *GenreRepository) DeleteGenre(id int) error {
	_, err := r.DB.Exec(`
		DELETE FROM Genre
		WHERE id = ?`, id,
	)
	return err
}
