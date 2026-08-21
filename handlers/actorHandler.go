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

type ActorHandler struct {
	Service *service.ActorService
}

type ActorResponse struct {
	Actor  models.Actor   `json:"actor"`
	Movies []models.Movie `json:"movies"`
}

func NewActorHandler(s *service.ActorService) *ActorHandler {
	return &ActorHandler{Service: s}
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
