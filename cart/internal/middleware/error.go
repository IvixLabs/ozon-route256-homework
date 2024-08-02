package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"route256/logger/pkg/logger"
)

type ErrorWrapper func(w http.ResponseWriter, req *http.Request) error

var (
	ErrValidation     = errors.New("validation error")
	ErrWrongArgument  = errors.New("wrong argument")
	ErrEntityNotFound = errors.New("entity is not found")
)

func (ew ErrorWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ew(w, r); err != nil {
		ctx := r.Context()

		switch {
		case errors.Is(err, ErrEntityNotFound):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrValidation):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, ErrWrongArgument):
			w.WriteHeader(http.StatusPreconditionFailed)
		default:
			logger.Errorw(ctx, "http error", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		type ErrorResponse struct {
			Message string `json:"message"`
		}

		bytesRes, jsonErr := json.Marshal(ErrorResponse{Message: err.Error()})
		if jsonErr != nil {
			log.Panic(jsonErr)
		}

		_, writeErr := w.Write(bytesRes)
		if writeErr != nil {
			log.Panic(writeErr)
		}
	}
}
