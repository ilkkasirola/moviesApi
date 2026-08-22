package repository

import (
	"database/sql"
	"movies-api/models"
)

type ActorRepository struct {
	DB *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{DB: db}
}

func (r *ActorRepository) GetAll() ([]models.Actor, error) {
	rows, err := r.DB.Query(`
		SELECT id, name, birthDate From Actor
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.Actor

	for rows.Next() {
		var a models.Actor
		err := rows.Scan(&a.ID, &a.Name, &a.BirthDate)
		if err != nil {
			return nil, err
		}
		actors = append(actors, a)
	}
	return actors, rows.Err()
}

func (r *ActorRepository) GetByID(id int) (models.Actor, error) {
	var a models.Actor

	err := r.DB.QueryRow(`
		SELECT id, name, birthDate FROM Actor WHERE id = ?
		`, id).Scan(&a.ID, &a.Name, &a.BirthDate)

	return a, err
}

func (r *ActorRepository) GetMoviesByActorID(actorID int) ([]models.Movie, error) {
	rows, err := r.DB.Query(`
		SELECT m.id, m.title, m.releaseYear, m.duration FROM Movie m
		INNER JOIN MovieActor ma ON m.id = ma.movieID WHERE ma.actorId = ?
		`, actorID)
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

func (r *ActorRepository) CreateActor(name string, birthDate any) (int64, error) {
	result, err := r.DB.Exec(
		"INSERT OR IGNORE INTO Actor (name, birthDate) VALUES (?, ?)",
		name, birthDate,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ActorRepository) GetActorID(name string, birthDate any) (int, error) {
	var id int
	err := r.DB.QueryRow(
		"SELECT id FROM Actor WHERE name = ? AND birthDate = ?",
		name, birthDate,
	).Scan(&id)
	return id, err
}

func (r *ActorRepository) LinkActor(movieID int64, actorID int) error {
	_, err := r.DB.Exec(
		"INSERT OR IGNORE INTO MovieActor (movieID, ActorID) VALUES (?, ?)",
		movieID, actorID,
	)
	return err
}
