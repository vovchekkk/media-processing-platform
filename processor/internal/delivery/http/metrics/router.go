package metrics

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
)

func RegisterRoutes(r chi.Router, log *slog.Logger) {
	r.Handle("/metrics", promhttp.Handler())
}
