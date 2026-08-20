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
