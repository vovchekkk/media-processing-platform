package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"media-processing-platform/server/internal/delivery/http/shared"
	"media-processing-platform/server/internal/dto"
	"media-processing-platform/server/internal/service"
)

type AuthHandler struct {
	log         *slog.Logger
	authService *service.AuthService
}

func NewAuthHandler(log *slog.Logger, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		log:         log,
		authService: authService,
	}
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает аккаунт пользователя с переданными DTO данными
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.User  true  "Данные для регистрации"
// @Success      200      "Пользователь успешно зарегистрирован"
// @Failure      400      {object}  map[string]string  "Неверный формат запроса или валидация не пройдена"
// @Failure      500      {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /register [post]
func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	req, ok := shared.DecodeAndValidate[dto.User](w, r)
	if !ok {
		return
	}

	if err := authHandler.authService.Register(r.Context(), &req); err != nil {
		authHandler.log.Error("registration failed", "error", err)
		shared.SendError(w, r, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Проверяет учетные данные и возвращает UUID токен сессии
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.User  true  "Данные для входа"
// @Success      200      {object}  map[string]string  "Возвращает token в JSON"
// @Failure      400      {object}  map[string]string  "Неверный формат запроса"
// @Failure      401      {object}  map[string]string  "Неверный логин или пароль"
// @Router       /login [post]
func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	req, ok := shared.DecodeAndValidate[dto.User](w, r)
	if !ok {
		return
	}

	token, err := authHandler.authService.Login(r.Context(), &req)
	if err != nil {
		authHandler.log.Error("unauthorized", "error", err)
		shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token.String()}); err != nil {
		authHandler.log.Error("failed to encode login response", "error", err)
	}
}

// Logout godoc
// @Summary      Выход из системы
// @Description  Завершает текущую сессию пользователя и инвалидирует UUID токен сессии
// @Tags         auth
// @Security     BearerAuth
// @Success      204  "Сессия успешно завершена"
// @Failure      401  {object} map[string]string "Пользователь не авторизован"
// @Failure      500  {object} map[string]string "Ошибка при завершении сессии"
// @Router       /logout [post]
func (authHandler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, ok := shared.GetSessionID(r.Context())
	if !ok {
		shared.SendError(w, r, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := authHandler.authService.Logout(r.Context(), sessionID); err != nil {
		authHandler.log.Error("failed to logout", "error", err)
		shared.SendError(w, r, http.StatusInternalServerError, "failed to logout")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
