package task

import (
	"log/slog"
	"net/http"
	"errors"
	"io"

	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"

	resp "media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/domain"
)

// New
// @Summary Create task
// @Description Creates a new media processing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Task creation request"
// @Success 200 {object} CreateResponse
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/tasks/ [post]
func New(log *slog.Logger, taskRepository repository.Task, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("empty request body", "error", err)
			
			render.JSON(w, r, resp.Error("empty request body"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body", "error", err)
			
			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		if err := validate.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("validation error", "error", validateErr)
			
			render.JSON(w, r, resp.Error("validation error"))

			return
		}

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
		respondOK(w, r, id)
	}
}

func respondOK(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	render.JSON(w, r, CreateResponse{
		ID: id,
		Response: resp.Success(),
	})
}