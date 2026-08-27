package auth

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"media-processing-platform/server/internal/delivery/http/shared"
	"media-processing-platform/server/internal/service"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			tokenStr = strings.TrimSpace(tokenStr)

			parsedToken, err := uuid.Parse(tokenStr)
			if err != nil {
				shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			session, err := authService.ValidateToken(r.Context(), parsedToken)
			if err != nil || session == nil {
				shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := r.Context()

			ctx = shared.SetUserID(ctx, session.UserID)
			ctx = shared.SetSessionID(ctx, session.ID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
