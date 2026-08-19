package repository

import (
	"database/sql"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

// Create movie

func (r *MovieRepository) Create(title string, year int, duration int) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT INTO Movie (title, releaseYear, duration) VALUES (?, ?, ?)",
		title, year, duration,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *MovieRepository) CreateActor(name string, birthDate any) error {
	_, err := r.DB.Exec(
		"INSERT OR IGNORE INTO Actor (name, birthDate) VALUES (?, ?)",
		name, birthDate,
	)
	return err
}

func (r *MovieRepository) GetActorID(name string, birthDate any) (int, error) {
	var id int
	err := r.DB.QueryRow(
		"SELECT id FROM Actor WHERE name = ? AND birthDate = ?",
		name, birthDate,
	).Scan(&id)
	return id, err
}

func (r *MovieRepository) LinkActor(movieID int64, actorID int) error {
	_, err := r.DB.Exec(
		"INSERT OR IGNORE INTO MovieActor (movieID, ActorID) VALUES (?, ?)",
		movieID, actorID,
	)
	return err
}
