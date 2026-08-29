package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	_ "media-processing-platform/server/docs"
	"media-processing-platform/server/internal/config"
	router "media-processing-platform/server/internal/delivery/http"
	"media-processing-platform/server/internal/service"
	postgresInfrastructure "media-processing-platform/server/internal/infrastructure/postgres"
	"media-processing-platform/server/internal/infrastructure/rabbitmq"
	postgresRepo "media-processing-platform/server/internal/repository/postgres"
	redisRepo "media-processing-platform/server/internal/repository/redis"
	redisInfrastructure "media-processing-platform/server/internal/infrastructure/redis"
)

const (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

// @title           Media Processing Platform API
// @version         1.0
// @description     API Server for Media Processing Platform

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <UUID_токен_сессии>
func main() {
	cfg := config.MustLoad()

	ctx := context.Background()

	logger := setupLogger(cfg.Env)
	logger = logger.With(slog.String("env", cfg.Env))

	logger.Info("initializing server", slog.String("address", cfg.HTTPServer.Address()))
	logger.Debug("logger debug mode enabled")

	postgresDB := postgresInfrastructure.InitDB(cfg.DatabaseConfig, logger)

	userRepository := postgresRepo.NewGormUserRepository(postgresDB)
	taskRepository := postgresRepo.NewGormTaskRepository(postgresDB)

	redisDB := redisInfrastructure.InitDB(ctx, cfg.RedisConfig, logger)

	sessionRepository := redisRepo.NewRedisSessionRepository(redisDB)

	connManager, err := rabbitmq.NewConnectionManager(cfg.RabbitMQConfig, logger)
	if err != nil {
		logger.Error("fatal: failed to initialize RabbitMQ connection manager", "error", err)
		os.Exit(1)
	}
	defer connManager.Close()

	taskProducer := rabbitmq.NewProducer(connManager, cfg.RabbitMQConfig.QueueName, logger)

	authService := service.NewAuthService(cfg.AuthConfig, userRepository, sessionRepository)
	taskService := service.NewTaskService(taskRepository, taskProducer)

	appRouter := router.InitRouter(logger, authService, taskService)

	logger.Info("starting server", slog.String("address", cfg.HTTPServer.Address()))

	if err := http.ListenAndServe(cfg.HTTPServer.Address(), appRouter); err != nil {
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
