package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
	httpSwagger "github.com/swaggo/http-swagger"

	"media-processing-platform/processor/internal/delivery/http/metrics"
)

func InitRouter(log *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(slogchi.New(log))

	router.Get("/swagger/*", httpSwagger.Handler())

	router.Group(func(r chi.Router) {
		metrics.RegisterRoutes(r, log)
	})

	return router
}
