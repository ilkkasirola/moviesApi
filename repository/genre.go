package repository

import (
	"encoding/json"
	"log"
	"movies-api/database"
	"movies-api/models"
	"net/http"
	"strconv"
	"strings"
)

func GetAllGenres(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := database.DB.Query("SELECT id, name FROM Genre")
	if err != nil {
		log.Println("Error querying genres", err)
		http.Error(w, "Error retrrieving genres", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			log.Println("Error scanning movie:", err)
			continue
		}
		genres = append(genres, genre)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("error during rows iteration", err)
	}
	json.NewEncoder(w).Encode(genres)
}

func GetGenreById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := r.URL.Path
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}
	var genre models.Genre
	err = database.DB.QueryRow("SELECT id, name FROM Genre WHERE id = ?", id).Scan(&genre.ID, &genre.Name)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(genre)
}
