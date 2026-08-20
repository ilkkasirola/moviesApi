package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes(mh *handlers.MovieHandler, ah *handlers.ActorHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", mh.GetAllMovies)
	mux.HandleFunc("GET /api/movies/{id}", mh.GetMovieByID)
	mux.HandleFunc("POST /api/movies", mh.CreateMovie)

	mux.HandleFunc("GET /api/actors", ah.GetAllActors)
	mux.HandleFunc("GET /api/actors/{id}", ah.GetActorByID)
	return mux
}
