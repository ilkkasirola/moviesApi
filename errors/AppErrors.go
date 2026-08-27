package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrValidation = &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "invalid request data",
		Status:  http.StatusBadRequest,
	}

	ErrMovieExists = &AppError{
		Code:    "MOVIE_EXISTS",
		Message: "movie already exists",
		Status:  http.StatusConflict,
	}

	ErrInvalidBody = &AppError{
		Code:    "INVALID_BODY",
		Message: "invalid request body",
		Status:  http.StatusBadRequest,
	}

	ErrInvalidDate = &AppError{
		Code:    "INVALID_DATE",
		Message: "birthdate must be YYYY-MM-DD",
		Status:  http.StatusBadRequest,
	}

	ErrIvalidID = &AppError{
		Code:    "INVALID_ID",
		Message: "movie id must be a positive integer",
		Status:  http.StatusNotFound,
	}

	ErrMovieNotFound = &AppError{
		Code:    "MOVIE_NOT_FOUND",
		Message: "movie not found",
		Status:  http.StatusNotFound,
	}

	ErrGenreNotFound = &AppError{
		Code:    "GENRE_NOT_FOUND",
		Message: "genre not found",
		Status:  http.StatusNotFound,
	}

	ErrGenreExists = &AppError{
		Code:    "GENRE_EXISTS",
		Message: "genre already exists",
		Status:  http.StatusBadRequest,
	}

	ErrActorExists = &AppError{
		Code:    "ACTOR_EXISTS",
		Message: "actor already exists",
		Status:  http.StatusBadRequest,
	}

	ErrActorNotFound = &AppError{
		Code:    "ACTOR_NOT_FOUND",
		Message: "actor not found",
		Status:  http.StatusNotFound,
	}
)

func countedAssociated(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("%d associated %s", n, noun)
	}
	return fmt.Sprintf("%d associated %ss", n, noun)
}

func NewGenreInUse(name string, movieCount int) *AppError {
	return &AppError{
		Code:    "GENRE_IN_USE",
		Message: fmt.Sprintf("Cannot delete genre '%s' because it has %s", name, countedAssociated(movieCount, "movie")),
		Status:  http.StatusBadRequest,
	}
}

func NewActorInUse(name string, movieCount int) *AppError {
	noun := "movies"
	if movieCount == 1 {
		noun = "movie"
	}
	return &AppError{
		Code:    "ACTOR_IN_USE",
		Message: fmt.Sprintf("Unable to delete actor '%s' as they are associated with %d %s", name, movieCount, noun),
		Status:  http.StatusBadRequest,
	}
}

func NewMovieInUse(title string, actorCount, genreCount int) *AppError {
	var parts []string
	if actorCount > 0 {
		parts = append(parts, countedAssociated(actorCount, "actor"))
	}
	if genreCount > 0 {
		parts = append(parts, countedAssociated(genreCount, "genre"))
	}
	return &AppError{
		Code:    "MOVIE_IN_USE",
		Message: fmt.Sprintf("Cannot delete movie '%s' because it has %s.", title, strings.Join(parts, " and ")),
		Status:  http.StatusBadRequest,
	}
}

func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if e, ok := err.(*AppError); ok {
		w.WriteHeader(e.Status)
		json.NewEncoder(w).Encode(e)
		return
	}

	w.WriteHeader(500)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "internal server error",
	})
}
