package task

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"

	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, log *slog.Logger, taskRepo repository.Task, validate *validator.Validate, taskProcessor *service.TaskProcessor) {
	r.Post("/task", NewCreateHandler(log, taskRepo, validate, taskProcessor))
	r.Get("/status/{task_id}", NewGetStatusHandler(log, taskRepo, validate))
	r.Get("/result/{task_id}", NewGetResultHandler(log, taskRepo, validate))
}
