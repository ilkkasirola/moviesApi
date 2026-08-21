package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes(mh *handlers.MovieHandler, ah *handlers.ActorHandler, gh *handlers.GenreHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", handlers.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", handlers.GetMovieById)
	mux.HandleFunc("POST /api/movies", handlers.CreateMovie)
	mux.HandleFunc("GET /api/genres", gh.GetAllGenres)

	mux.HandleFunc("GET /api/actors", ah.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", ah.GetActorByID)
	mux.HandleFunc("GET /api/genres/{id}", gh.GetGenreById)
	mux.HandleFunc("POST /api/genres", gh.CreateGenre)
	return mux
}
