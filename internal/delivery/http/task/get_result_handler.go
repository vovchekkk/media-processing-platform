package task

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	resp "media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/repository"
)

// NewGetResultHandler
// @Summary Get task result
// @Description Returns the processing result of a task by its ID
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task UUID" format(uuid)
// @Success 200 {object} GetResultResponse
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/tasks/result/{task_id}/ [get]
func NewGetResultHandler(log *slog.Logger, taskRepository repository.Task, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := resp.BindPathUUID(w, r, log, "task_id")
		if !ok {
			return
		}

		result, err := taskRepository.GetTaskResultByID(id)
		if err != nil {
			log.Error("failed to get task result", "error", err)

			render.JSON(w, r, resp.Error("failed to get task result"))

			return
		}

		log.Info("task result successfully retrieved", "task_id", id, "result", result)
		getResultRespondOK(w, r, result)
	}
}

func getResultRespondOK(w http.ResponseWriter, r *http.Request, result string) {
	render.JSON(w, r, GetResultResponse{
		Result:   result,
		Response: resp.Success(),
	})
}
