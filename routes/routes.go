package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", handlers.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", handlers.GetMovieById)
	// mux.HandleFunc("GET /api/actors", GetAllActors)

	mux.HandleFunc("POST /api/movies", handlers.CreateMovie)

	return mux
}
