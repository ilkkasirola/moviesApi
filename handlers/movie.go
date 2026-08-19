package handlers

import (
	"encoding/json"
	appErr "movies-api/errors"
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

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if e, ok := err.(*appErr.AppError); ok {
		w.WriteHeader(e.Status)
		json.NewEncoder(w).Encode(e)
		return
	}

	w.WriteHeader(500)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		WriteError(w, appErr.ErrInvalidBody)
		return
	}

	movieID, err := h.Service.CreateMovie(movie)
	if err != nil {
		WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": movieID,
	})
}

// 	movie.ID = int(movieID)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(movie)
// }
