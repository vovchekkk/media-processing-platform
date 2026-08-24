package task

import (
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	resp "media-processing-platform/internal/delivery/http/shared"
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
// @Failure 401 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Security BearerAuth
// @Router /status/{task_id} [get]
func NewGetResultHandler(log *slog.Logger, taskRepository repository.Task, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := resp.BindPathUUID(w, r, "task_id")
		if !ok {
			return
		}

		result, err := taskRepository.GetTaskResultByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, resp.Error("task not found"))
				return
			}

			log.Error("failed to get task result", "error", err)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get task result"))
			return
		}

		log.Info("task result successfully retrieved", "task_id", id, "result", result)
		getResultRespondOK(w, r, result)
	}
}

func getResultRespondOK(w http.ResponseWriter, r *http.Request, result string) {
	render.JSON(w, r, GetResultResponse{
		Result: result,
	})
}
