package metrics

import (
	"log/slog"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RegisterRoutes(r chi.Router, log *slog.Logger) {
	r.Handle("/metrics", promhttp.Handler())
}
