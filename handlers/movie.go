package handlers

import (
	"encoding/json"
	"movies-api/models"
	"movies-api/service"
	"net/http"
)

type MovieHandler struct {
	Service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{Service: s}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	movieID, err := h.Service.CreateMovie(movie)
	if err != nil {
		if err == service.ErrMovieExists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "movie already exists",
				"code":  "MOVIE_EXISTS",
			})
			return
		}

		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	movie.ID = int(movieID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}
