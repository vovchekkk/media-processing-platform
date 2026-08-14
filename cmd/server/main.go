package main

import (
	"log/slog"
	"os"

	"media-processing-platform/internal/config"
	"media-processing-platform/internal/infrastructure/database"
	"media-processing-platform/internal/repository"
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

	db := database.InitDB(logger)
	taskRepository := repository.NewGormTaskRepository(db)

	_ = taskRepository // Use the taskRepository as needed
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