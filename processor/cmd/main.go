package main

import (
	"context"
	"log/slog"
	"media-processing-platform/processor/internal/config"
	database "media-processing-platform/processor/internal/infrastructure/postgres"
	"media-processing-platform/processor/internal/infrastructure/rabbitmq"
	"media-processing-platform/processor/internal/repository/postgres"
	"media-processing-platform/processor/internal/service"
	"os"
	"os/signal"
	"syscall"
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

	db := database.InitDB(cfg.DatabaseConfig, logger)

	taskRepository := postgres.NewGormTaskRepository(db)

	connManager, err := rabbitmq.NewConnectionManager(cfg.RabbitMQConfig, logger)
	if err != nil {
		logger.Error("fatal: failed to initialize RabbitMQ connection manager", "error", err)
		os.Exit(1)
	}
	defer connManager.Close()

	taskConsumer := rabbitmq.NewConsumer(connManager, cfg.RabbitMQConfig.QueueName, logger)

	processorService := service.NewProcessorService(taskRepository, logger)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := taskConsumer.StartListening(ctx, processorService.Process)
		if err != nil {
			logger.Error("consumer crashed", "error", err)
		}
	}()

	logger.Info("processor worker started successfully")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("processor worker shutting down")
	cancel()
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
