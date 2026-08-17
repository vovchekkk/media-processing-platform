package task

import (
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	resp "media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
)

// NewGetStatusHandler
// @Summary Get task status
// @Description Returns the current processing status of a task by its ID
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task UUID" format(uuid)
// @Success 200 {object} GetStatusResponse
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /status/{task_id} [get]
func NewGetStatusHandler(log *slog.Logger, taskRepository repository.Task, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := resp.BindPathUUID(w, r, log, "task_id")
		if !ok {
			return
		}

		result, err := taskRepository.GetTaskStatusByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("task not found"))
				return
			}

			log.Error("failed to get task status", "error", err)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get task status"))
			return
		}

		log.Info("task status successfully retrieved", "task_id", id, "status", result)
		getStatusRespondOK(w, r, result)
	}
}

func getStatusRespondOK(w http.ResponseWriter, r *http.Request, status domain.TaskStatus) {
	render.JSON(w, r, GetStatusResponse{
		TaskStatus: status,
	})
}
