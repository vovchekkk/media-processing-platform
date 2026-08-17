package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-playground/validator"

	_ "media-processing-platform/docs"
	"media-processing-platform/internal/config"
	router "media-processing-platform/internal/delivery/http"
	"media-processing-platform/internal/infrastructure/postgres"
	"media-processing-platform/internal/repository/postgres"
	"media-processing-platform/internal/service"
)

const (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger = logger.With(slog.String("env", cfg.Env))

	logger.Info("initializing server", slog.String("address", cfg.Address()))
	logger.Debug("logger debug mode enabled")

	db := database.InitDB(cfg.DatabaseConfig, logger)

	userRepository := postgres.NewGormUserRepository(db)
	taskRepository := postgres.NewGormTaskRepository(db)
	sessionRepository := postgres.NewGormSessionRepository(db)

	authService := service.NewAuthService(userRepository, sessionRepository)

	validate := validator.New()

	taskProcessor := service.NewTaskProcessor(cfg.TaskProcessorConfig, taskRepository, logger)

	appRouter := router.InitRouter(logger, authService, taskRepository, validate, taskProcessor)

	logger.Info("starting server", slog.String("address", cfg.Address()))

	if err := http.ListenAndServe(cfg.Address(), appRouter); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		slog.Error("unknown environment", "environment", env)
		return nil
	}

	return logger
}
