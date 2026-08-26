package handlers

import (
	"encoding/json"
	"errors"
	appErr "movies-api/errors"
	"movies-api/models"
	"movies-api/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type GenreHandler struct {
	Service  *service.GenreService
	Validate *validator.Validate
}

type GenreResponse struct {
	Genre models.Genre `json:"Genre"`
}

func NewGenreHandler(s *service.GenreService, v *validator.Validate) *GenreHandler {
	return &GenreHandler{Service: s, Validate: v}
}

func (h *GenreHandler) GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.Service.GetAllGenres()
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetGenreById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	genre, err := h.Service.GetGenreByID(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(GenreResponse{Genre: genre})
}

func (h *GenreHandler) GetGenreByMovieId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	genres, err := h.Service.GetGenreByMovieID(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var genre models.Genre
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	genreID, err := h.Service.CreateGenre(genre)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": genreID,
	})
}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	err = h.Service.DeleteGenre(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) PatchGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	var patch models.GenrePatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	if patch.Name == nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	if strings.TrimSpace(*patch.Name) == "" {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	genre, err := h.Service.PatchGenre(r.Context(), id, *patch.Name)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrGenreNotFound):
			appErr.WriteError(w, appErr.ErrGenreNotFound)
		case errors.Is(err, appErr.ErrGenreExists):
			appErr.WriteError(w, appErr.ErrGenreExists)
		default:
			appErr.WriteError(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}
