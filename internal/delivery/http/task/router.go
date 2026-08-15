package task

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"

	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/service"
)

func InitRouter(log *slog.Logger, taskRepo repository.Task, validate *validator.Validate, taskProcessor *service.TaskProcessor) http.Handler {
	router := chi.NewRouter()

	router.Post("/task", NewCreateHandler(log, taskRepo, validate, taskProcessor))
	router.Get("/status/{task_id}", NewGetStatusHandler(log, taskRepo, validate))
	router.Get("/result/{task_id}", NewGetResultHandler(log, taskRepo, validate))

	return router
}