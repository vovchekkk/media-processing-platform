package http

import (
	"net/http"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/samber/slog-chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	httpSwagger "github.com/swaggo/http-swagger"

	"media-processing-platform/internal/repository"
	"media-processing-platform/internal/delivery/http/task"
)

func InitRouter(log *slog.Logger, taskRepo repository.Task, validate *validator.Validate) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(slogchi.New(log))

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Route("/api", func (r chi.Router) {
		r.Mount("/tasks", task.InitRouter(log, taskRepo, validate))
	})

	return router
}