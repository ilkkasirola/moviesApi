package handlers

import (
	"encoding/json"
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
