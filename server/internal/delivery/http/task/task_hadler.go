package task

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"media-processing-platform/server/internal/delivery/http/shared"
	"media-processing-platform/server/internal/domain"
	"media-processing-platform/server/internal/dto"
	"media-processing-platform/server/internal/service"
)

type TaskHandler struct {
	log         *slog.Logger
	taskService *service.TaskService
}

func NewTaskHandler(log *slog.Logger, taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		log:         log,
		taskService: taskService,
	}
}

// Create godoc
// @Summary Create task
// @Description Creates a new image processing task. The image must be provided as a Base64-encoded string.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.Task true "Task creation request"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /task [post]
func (taskHandler *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	req, ok := shared.DecodeAndValidate[dto.Task](w, r)
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(r.Context())
	if !ok {
		shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := taskHandler.taskService.Create(r.Context(), &req, userID)
	if err != nil {
		taskHandler.log.Error("failed to create task", "error", err)
		shared.SendError(w, r, http.StatusInternalServerError, "failed to create task")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{"task_id": id.String()}); err != nil {
		taskHandler.log.Error("failed to encode create response", "error", err)
	}
}

// GetResult godoc
// @Summary Get task result
// @Description Returns the result of an image processing task by its ID
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task UUID" format(uuid)
// @Success 200 {object} GetResultResponse
// @Failure 400 {object} map[string]string "Invalid task ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /result/{task_id} [get]
func (taskHandler *TaskHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	id, ok := shared.BindPathUUID(w, r, "task_id")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(r.Context())
	if !ok {
		shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := taskHandler.taskService.GetResult(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			taskHandler.log.Error("task not found", "error", err)
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
// @Success 200 {object} GetStatusResponse
// @Failure 400 {object} map[string]string "Invalid task ID"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /status/{task_id} [get]
func (taskHandler *TaskHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := shared.BindPathUUID(w, r, "task_id")
	if !ok {
		return
	}

	userID, ok := shared.GetUserID(r.Context())
	if !ok {
		shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := taskHandler.taskService.GetStatus(r.Context(), id, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTaskNotFound):
			taskHandler.log.Error("task not found", "error", err)
			shared.SendError(w, r, http.StatusNotFound, "task not found")
		default:
			taskHandler.log.Error("failed to get task status", "error", err)
			shared.SendError(w, r, http.StatusInternalServerError, "failed to get task status")
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"status": string(status)}); err != nil {
		taskHandler.log.Error("failed to encode get_status response", "error", err)
	}
}
