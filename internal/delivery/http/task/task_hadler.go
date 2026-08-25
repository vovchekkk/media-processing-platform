package task

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/service"
)

type TaskHandler struct {
	log         *slog.Logger
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// Create godoc
// @Summary Create task
// @Description Creates a new media processing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Task creation request"
// @Success 201 {object} CreateResponse
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Security BearerAuth
// @Router /task [post]
func (taskHandler *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	id, err := taskHandler.taskService.Create(r.Context())
	if err != nil {
		taskHandler.log.Error("failed to create task", "error", err)
		shared.SendError(w, r, http.StatusInternalServerError, "failed to create task")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{"id": id.String()}); err != nil {
		taskHandler.log.Error("failed to encode create response", "error", err)
	}
}

// GetResult godoc
// @Summary Get task result
// @Description Returns the processing result of a task by its ID
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task UUID" format(uuid)
// @Success 200 {object} GetResultResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Security BearerAuth
// @Router /result/{task_id} [get]
func (taskHandler *TaskHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	id, ok := shared.BindPathUUID(w, r, "task_id")
	if !ok {
		return
	}

	result, err := taskHandler.taskService.GetResult(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			shared.SendError(w, r, http.StatusNotFound, "task not found")
		default:
			taskHandler.log.Error("failed to get task result", "error", err)
			shared.SendError(w, r, http.StatusInternalServerError, "failed to get task result")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"result": result}); err != nil {
		taskHandler.log.Error("failed to encode get_result response", "error", err)
	}
}

// GetStatus godoc
// @Summary Get task status
// @Description Returns the current processing status of a task by its ID
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task UUID" format(uuid)
// @Success 200 {object} GetResultResponse
// @Failure 400 {object} resp.Response
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Security BearerAuth
// @Router /status/{task_id} [get]
func (taskHandler *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := shared.BindPathUUID(w, r, "task_id")
	if !ok {
		return
	}

	status, err := taskHandler.taskService.GetStatus(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			shared.SendError(w, r, http.StatusNotFound, "task not found")
		default:
			taskHandler.log.Error("failed to get task status", "error", err)
			shared.SendError(w, r, http.StatusInternalServerError, "failed to get task status")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"status": string(status)}); err != nil {
		taskHandler.log.Error("failed to encode get_status response", "error", err)
	}
}
