package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	"github.com/samber/slog-chi"
	httpSwagger "github.com/swaggo/http-swagger"

	"media-processing-platform/internal/delivery/http/task"
	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/service"
)

func InitRouter(log *slog.Logger, taskRepo repository.Task, validate *validator.Validate, taskProcessor *service.TaskProcessor) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(slogchi.New(log))

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/", func (r chi.Router) {
		r.Mount("/", task.InitRouter(log, taskRepo, validate, taskProcessor))
	})

	return router
}