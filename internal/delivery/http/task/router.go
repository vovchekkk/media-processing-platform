package task

import (
	"github.com/go-chi/chi/v5"
	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, taskService *service.TaskService) {
	r.Post("/task", NewTaskHandler(taskService).Create)
	r.Get("/status/{task_id}", NewTaskHandler(taskService).GetStatus)
	r.Get("/result/{task_id}", NewTaskHandler(taskService).GetResult)
}
