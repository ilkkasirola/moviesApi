package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes(mh *handlers.MovieHandler, ah *handlers.ActorHandler, gh *handlers.GenreHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", mh.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", mh.GetMovieByID)
	mux.HandleFunc("POST /api/movies", mh.CreateMovie)

	mux.HandleFunc("GET /api/actors", ah.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", ah.GetActorByID)
	mux.HandleFunc("POST /api/actors", ah.CreateActor)

	mux.HandleFunc("GET /api/genres", gh.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", gh.GetGenreById)
	mux.HandleFunc("POST /api/genres", gh.CreateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", gh.DeleteGenre)
	mux.HandleFunc("PATCH /api/genres/{id}", gh.PatchGenre)
	return mux
}
