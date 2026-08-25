package auth

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, log *slog.Logger, authService *service.AuthService) {
	handler := NewAuthHandler(log, authService)
	
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
}
