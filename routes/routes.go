package routes

import (
	"movies-api/handlers"
	"net/http"
)

func SetupRoutes(mh *handlers.MovieHandler, ah *handlers.ActorHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/movies", mh.G)
	// mux.HandleFunc("GET /api/movies/{id}", handlers.GetMovieById)
	mux.HandleFunc("POST /api/movies", mh.CreateMovie)

	// mux.HandleFunc("PATCH /api/movies/{id}", handlers.PatchMovie)
	return mux
}
