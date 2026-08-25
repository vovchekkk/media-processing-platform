package task

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, log *slog.Logger, taskService *service.TaskService) {
	r.Post("/task", NewTaskHandler(log, taskService).Create)
	r.Get("/status/{task_id}", NewTaskHandler(log, taskService).GetStatus)
	r.Get("/result/{task_id}", NewTaskHandler(log, taskService).GetResult)
}
