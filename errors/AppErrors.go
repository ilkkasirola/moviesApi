package errors

import (
	"encoding/json"
	"net/http"
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
