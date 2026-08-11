package handlers

import (
	"encoding/json"
	"log"
	"movies-api/database"
	"movies-api/models"
	"net/http"
	"strconv"
	"strings"
)

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT id, title, releaseYear, duration FROM Movie")
	if err != nil {
		log.Println("Error querying movies", err)
		http.Error(w, "Error retrrieving movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Duration); err != nil {
			log.Println("Error scanning movie:", err)
			continue
		}

		// fetch genres
		genres, err := getMovieGenres(movie.ID)
		if err != nil {
			log.Println("Error fetching genres:", err)
		}
		movie.Genres = genres

		actors, err := getMovieActors(movie.ID)
		if err != nil {
			log.Println("Error fetching actors:", err)
		}
		movie.Actors = actors

		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error during rows iteration", err)
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovieGenres(movieID int) ([]models.Genre, error) {
	query := `
	SELECT g.id, g.name FROM Genre g
	INNER JOIN MovieGenre mg On g.id = mg.GenreId
	WHERE mg.movieId = ?
	`
	rows, err := database.DB.Query(query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error during rows iteration", err)
	}
	return genres, nil
}

// func GetAllActors(w http.ResponseWriter, r *http.Request) {
// 	query := `
// 	SELECT id, name, birthDate FROM Actor
// 	INNER JOIN
// 	`
// }

func getMovieActors(movieID int) ([]models.Actor, error) {
	query := `
	SELECT a.id, a.name, a.birthDate FROM Actor a
	INNER JOIN MovieActor ma ON a.id = ma.actorId
	WHERE ma.movieId = ?
	`
	rows, err := database.DB.Query(query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.Actor
	for rows.Next() {
		var actor models.Actor
		if err := rows.Scan(&actor.ID, &actor.Name, &actor.BirthDate); err != nil {
			return nil, err
		}
		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error during rows iteration", err)
	}
	return actors, nil
}

func GetMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}
	var movie models.Movie
	err = database.DB.QueryRow("SELECT id, title, releaseYear, duration FROM movie WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Duration)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	genres, _ := getMovieGenres(movie.ID)
	movie.Genres = genres
	actors, _ := getMovieActors(movie.ID)
	movie.Actors = actors

	json.NewEncoder(w).Encode(movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var movie models.Movie
	var actor models.Actor
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := database.DB.Exec(
		"INSERT INTO Movie (title, releaseYear, duration) VALUES (?, ?, ?)",
		movie.Title, movie.ReleaseYear, movie.Duration,
	)
	if err != nil {
		log.Println("Error inserting  movie", err)
		http.Error(w, "Error creating movie", http.StatusInternalServerError)
		return
	}
	_, err = database.DB.Exec(
		"INSERT INTO Actor (name, birthDate) VALUES (?, ?)",
		actor.Name, actor.BirthDate,
	)
	if err != nil {
		log.Println("Error inserting actor", err)
		http.Error(w, "Error creating actor", http.StatusInternalServerError)
		return
	}

	lastId, _ := result.LastInsertId()
	movie.ID = int(lastId)
	movie.Genres = []models.Genre{}
	movie.Actors = []models.Actor{}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}
