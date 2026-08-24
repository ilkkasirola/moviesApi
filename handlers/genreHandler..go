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
)

type GenreHandler struct {
	Service *service.GenreService
}

type GenreResponse struct {
	Genre models.Genre `json:"Genre"`
}

func NewGenreHandler(s *service.GenreService) *GenreHandler {
	return &GenreHandler{Service: s}
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
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
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
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
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
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	var patch models.GenrePatch
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if patch.Name == nil {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(*patch.Name) == "" {
		http.Error(w, "name canno be empty", http.StatusBadRequest)
		return
	}

	genre, err := h.Service.PatchGenre(r.Context(), id, *patch.Name)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrGenreNotFound):
			http.Error(w, "genre not found", http.StatusNotFound)
		case errors.Is(err, appErr.ErrGenreExists):
			http.Error(w, "genre name already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}
