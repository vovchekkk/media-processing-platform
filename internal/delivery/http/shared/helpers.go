package shared

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

func BindAndValidate(w http.ResponseWriter, r *http.Request, log *slog.Logger, target any, validate *validator.Validate) bool {
	err := render.DecodeJSON(r.Body, target)
	
	if errors.Is(err, io.EOF) {
		log.Error("empty request body", slog.Any("error", err))
		render.JSON(w, r, Error("empty request body"))
		return false
	}
	
	if err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		render.JSON(w, r, Error("failed to decode request body"))
		return false
	}

	if err := validate.Struct(target); err != nil {
		log.Error("validation error", slog.Any("error", err.Error()))
		render.JSON(w, r, Error("validation error"))
		return false
	}

	return true
}

func BindPathUUID(
	w http.ResponseWriter,
	r *http.Request,
	log *slog.Logger,
	param string,
) (uuid.UUID, bool) {
	value := chi.URLParam(r, param)

	id, err := uuid.Parse(value)
	if err != nil {
		log.Error("invalid UUID path parameter",
			"parameter", param,
			"value", value,
			"error", err,
		)

		render.JSON(w, r, Error("invalid task id"))
		return uuid.Nil, false
	}

	return id, true
}