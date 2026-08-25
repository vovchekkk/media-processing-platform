package auth

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, log *slog.Logger, authService *service.AuthService) {
	r.Post("/register", NewAuthHandler(log, authService).Register)
	r.Post("/login", NewAuthHandler(log, authService).Login)
}
