package auth

import (
	"github.com/go-chi/chi/v5"

	"media-processing-platform/internal/service"
)

func RegisterRoutes(r chi.Router, authService *service.AuthService) {
	r.Post("/register", NewAuthHandler(authService).Register)
	r.Post("/login", NewAuthHandler(authService).Login)
}
