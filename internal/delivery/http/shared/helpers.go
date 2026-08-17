package shared

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

var validate = validator.New()

func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var target T
	
	err := render.DecodeJSON(r.Body, &target)
	if err != nil {
		if errors.Is(err, io.EOF) {
			SendError(w, r, http.StatusBadRequest, "empty request body")
			return target, false
		}
		
		SendError(w, r, http.StatusBadRequest, "failed to decode request body")
		return target, false
	}

	if err := validate.Struct(target); err != nil {
		SendError(w, r, http.StatusBadRequest, "validation error")
		return target, false
	}

	return target, true
}

func BindPathUUID(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	value := chi.URLParam(r, param)

	id, err := uuid.Parse(value)
	if err != nil {
		SendError(w, r, http.StatusBadRequest, "invalid UUID path parameter")
		return uuid.Nil, false
	}

	return id, true
}

func SendError(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	render.Status(r, statusCode)
	render.JSON(w, r, Error(msg))
}