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

type ActorHandler struct {
	Service  *service.ActorService
	Validate *validator.Validate
}

type ActorResponse struct {
	Actor  models.Actor   `json:"actor"`
	Movies []models.Movie `json:"movies"`
}

func NewActorHandler(s *service.ActorService, v *validator.Validate) *ActorHandler {
	return &ActorHandler{Service: s, Validate: v}
}

func (h *ActorHandler) GetAllActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.Service.GetAllActors()
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetActorByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	actor, movies, err := h.Service.GetActorByID(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(ActorResponse{
		Actor:  actor,
		Movies: movies,
	})
}
func (h *ActorHandler) CreateActor(w http.ResponseWriter, r *http.Request) {
	var actor models.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}
	err = h.Validate.Struct(actor)
	if err != nil {
		appErr.WriteError(w, appErr.ErrValidation)
		return
	}

	genreID, err := h.Service.CreateActor(actor)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": genreID,
	})
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}
	err = h.Service.DeleteActor(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	var patch models.ActorPatch
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

	actor, err := h.Service.PatchActor(r.Context(), id, *patch.Name, *patch.BirthDate)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrGenreNotFound):
			http.Error(w, "actor not found", http.StatusNotFound)
		case errors.Is(err, appErr.ErrGenreExists):
			http.Error(w, "actor name already exists", http.StatusConflict)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}
