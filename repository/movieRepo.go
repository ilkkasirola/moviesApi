package repository

import (
	"database/sql"
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
