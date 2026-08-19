package main

import (
	"database/sql"
	"log"
	"net/http"
)

type apiConfig struct {
	db *sql.DB
}

func RunServer(db *sql.DB) {
	cfg := &apiConfig{db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/movies", cfg.MoviesPostHandler)
	mux.HandleFunc("GET /api/movies", cfg.MoviesGetHandler)
	mux.HandleFunc("GET /api/movies/", cfg.MoviesIdGetHandler)
	mux.HandleFunc("POST /api/genres", cfg.GenresPostHandler)
	mux.HandleFunc("GET /api/genres", cfg.GenresGetHandler)
	mux.HandleFunc("GET /api/genres/", cfg.GenresIdGetHandler)
	mux.HandleFunc("POST /api/actors", cfg.ActorsPostHandler)
	mux.HandleFunc("GET /api/actors", cfg.ActorsGetHandler)
	mux.HandleFunc("GET /api/actors/", cfg.ActorsIdGetHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func (cfg *apiConfig) MoviesPostHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) MoviesGetHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.Query(`SELECT * FROM Movies`)
	if err != nil {
		http.Error(w, "failed to query movies", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()
}
func (cfg *apiConfig) MoviesIdGetHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) GenresPostHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) GenresGetHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) GenresIdGetHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) ActorsPostHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) ActorsGetHandler(w http.ResponseWriter, r *http.Request) {

}
func (cfg *apiConfig) ActorsIdGetHandler(w http.ResponseWriter, r *http.Request) {

}
