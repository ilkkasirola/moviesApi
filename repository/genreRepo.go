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

func (r *GenreRepository) LinkGenre(movieID int64, genreID int) error {
	_, err := r.DB.Exec("INSERT OR IGNORE INTO MovieGenre (movieID, genreID) VALUES (?, ?)", movieID, genreID)
	return err
}
