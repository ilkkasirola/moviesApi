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
	name := r.URL.Query().Get("name")

	if name != "" {
		h.SearchActorByName(w, r)
		return
	}

	actors, err := h.Service.GetAllActors()
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetActorByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
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

	actorID, err := h.Service.CreateActor(actor)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{
		"id": actorID,
	})
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}
	force := r.URL.Query().Get("force") == "true"
	err = h.Service.DeleteActor(id, force)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ActorHandler) GetActorsByMovieID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	actors, err := h.Service.GetActorsByMovieID(id)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) SearchActorByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	actors, err := h.Service.SearchActorByName(name)
	if err != nil {
		appErr.WriteError(w, err)
		return
	}

	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) PatchActor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}

	var patch models.ActorPatch
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

	actor, err := h.Service.PatchActor(r.Context(), id, *patch.Name, patch.BirthDate)
	if err != nil {
		switch {
		case errors.Is(err, appErr.ErrActorNotFound):
			appErr.WriteError(w, appErr.ErrActorNotFound)
		case errors.Is(err, appErr.ErrActorExists):
			appErr.WriteError(w, appErr.ErrActorExists)
		default:
			appErr.WriteError(w, err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(actor); err != nil {
		appErr.WriteError(w, appErr.ErrInvalidBody)
		return
	}
}
