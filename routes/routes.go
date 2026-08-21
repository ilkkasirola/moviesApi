package routes

import (
	"movies-api/handlers"
	"movies-api/repository"
	"net/http"
)

func SetupRoutes(h *handlers.MovieHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", handlers.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", handlers.GetMovieById)
	mux.HandleFunc("POST /api/movies", handlers.CreateMovie)
	mux.HandleFunc("GET /api/genres", repository.GetAllGenres)

	// mux.HandleFunc("PATCH /api/movies/{id}", handlers.PatchMovie)
	return mux
}
