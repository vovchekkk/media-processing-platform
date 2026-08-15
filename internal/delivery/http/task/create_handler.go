package task

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"

	resp "media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/service"
)

// NewCreateHandler
// @Summary Create task
// @Description Creates a new media processing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Task creation request"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /task [post]
func NewCreateHandler(log *slog.Logger, taskRepository repository.Task, validate *validator.Validate, taskProcessor *service.TaskProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Error("failed to generate uuid", "error", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log = log.With("creating task", "task_id", id)

		task := domain.Task{
			ID: id,
		}

		if err := taskRepository.CreateTask(&task); err != nil {
			log.Error("failed to create task", "error", err)

			render.JSON(w, r, resp.Error("failed to create task"))

			return
		}

		log.Info("task successfully created and saved")

		taskProcessor.Process(id)

		createRespondOK(w, r, id)
	}
}

func createRespondOK(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, CreateResponse{
		ID: id,
	})
}
