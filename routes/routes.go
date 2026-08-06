package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(" /api/movies/", handleMovies)

	return mux
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/movies" {
		if r.Method == http.MethodGet {
			handlers.GetAllMovies(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateMovie(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if r.Method == http.MethodGet {
		handlers.GetMovieById(w, r)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
