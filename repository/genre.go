package repository

import (
	"encoding/json"
	"log"
	"movies-api/database"
	"movies-api/models"
	"net/http"
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
