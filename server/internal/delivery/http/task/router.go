package task

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"media-processing-platform/server/internal/service"
)

func RegisterRoutes(r chi.Router, log *slog.Logger, taskService *service.TaskService) {
	handler := NewTaskHandler(log, taskService)
	
	r.Post("/task", handler.Create)
	r.Get("/status/{task_id}", handler.GetStatus)
	r.Get("/result/{task_id}", handler.GetResult)
}
