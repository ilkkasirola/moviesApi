package handlers

import (
	"encoding/json"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/service"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	Service *service.MovieService
}

func NewMovieHandler(s *service.MovieService) *MovieHandler {
	return &MovieHandler{Service: s}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	movieID, err := h.Service.CreateMovie(movie)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": movieID,
	})
}

func (h *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.Service.GetAllMovies()
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		appErr.WriteError(w, appErr.ErrIvalidID)
		return
	}

	movie, err := h.Service.GetMovieByID(id)
	if err != nil {
		appErr.WriteError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// 	movie.ID = int(movieID)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(movie)
// }
