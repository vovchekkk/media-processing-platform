package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
	httpSwagger "github.com/swaggo/http-swagger"

	"media-processing-platform/server/internal/delivery/http/task"
	"media-processing-platform/server/internal/service"
	"media-processing-platform/server/internal/delivery/http/auth"
)

func InitRouter(log *slog.Logger, authService *service.AuthService, taskService *service.TaskService) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(slogchi.New(log))

	router.Get("/swagger/*", httpSwagger.Handler())

	auth.RegisterRoutes(router, log, authService)

	router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(authService))

		task.RegisterRoutes(r, log, taskService)
	})

	return router
}
