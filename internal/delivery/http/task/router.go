package task

import (
	"net/http"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"

	"media-processing-platform/internal/repository"
)

func InitRouter(log *slog.Logger, taskRepo repository.Task, validate *validator.Validate) http.Handler {
	router := chi.NewRouter()

	router.Post("/", New(log, taskRepo, validate))

	return router
}