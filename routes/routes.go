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
	mux.HandleFunc("GET /api/movies/{id}/actors", ah.GetActorsByMovieID)
	mux.HandleFunc("DELETE /api/movies/{id}", mh.DeleteMovie)
	mux.HandleFunc("PATCH /api/movies/{id}", mh.PatchMovie)

	mux.HandleFunc("GET /api/actors", ah.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", ah.GetActorByID)
	mux.HandleFunc("POST /api/actors", ah.CreateActor)
	mux.HandleFunc("DELETE /api/actors/{id}", ah.DeleteActor)
	mux.HandleFunc("PATCH /api/actors/{id}", ah.PatchActor)

	mux.HandleFunc("GET /api/genres", gh.GetAllGenres)
	mux.HandleFunc("GET /api/genres/{id}", gh.GetGenreById)
	mux.HandleFunc("POST /api/genres", gh.CreateGenre)
	mux.HandleFunc("DELETE /api/genres/{id}", gh.DeleteGenre)
	mux.HandleFunc("PATCH /api/genres/{id}", gh.PatchGenre)
	mux.HandleFunc("GET /api/movies/{id}/genres", gh.GetGenreByMovieId)

	return mux
}
