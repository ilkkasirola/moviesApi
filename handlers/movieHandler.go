package handlers

import (
	"encoding/json"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/service"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type MovieHandler struct {
	Service  *service.MovieService
	Validate *validator.Validate
}

func NewMovieHandler(s *service.MovieService, v *validator.Validate) *MovieHandler {
	return &MovieHandler{
		Service:  s,
		Validate: v,
	}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")

	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}
	err = h.Validate.Struct(movie)
	if err != nil {
		appErr.WriteError(w, appErr.ErrValidation)
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// 	movie.ID = int(movieID)

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(movie)
// }
