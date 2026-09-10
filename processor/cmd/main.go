package main

import (
	"context"
	"log/slog"
	"media-processing-platform/processor/internal/config"
	router "media-processing-platform/processor/internal/delivery/http"
	"media-processing-platform/processor/internal/filter"
	postgresInfrastructure "media-processing-platform/processor/internal/infrastructure/postgres"
	"media-processing-platform/processor/internal/infrastructure/rabbitmq"
	"media-processing-platform/processor/internal/metrics"
	postgresRepo "media-processing-platform/processor/internal/repository/postgres"
	"media-processing-platform/processor/internal/service"
	"net/http"
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

	logger.Info("initializing processor", slog.String("address", cfg.HTTPServer.Address()))
	logger.Debug("logger debug mode enabled")

	db := postgresInfrastructure.InitDB(cfg.DatabaseConfig, logger)

	taskRepository := postgresRepo.NewGormTaskRepository(db)

	connManager, err := rabbitmq.NewConnectionManager(cfg.RabbitMQConfig, logger)
	if err != nil {
		logger.Error("fatal: failed to initialize RabbitMQ connection manager", "error", err)
		os.Exit(1)
	}
	defer connManager.Close()

	taskConsumer := rabbitmq.NewConsumer(connManager, cfg.RabbitMQConfig.QueueName, logger)

	registry := filter.NewRegistry(
		&filter.Blur{},
		&filter.Negative{},
		&filter.FlipX{},
		&filter.Sharpen{},
	)

	imageService := service.NewImageService(registry)

	metricsInstance := metrics.New()

	processorService := service.NewProcessorService(taskRepository, imageService, metricsInstance, logger)

	appRouter := router.InitRouter(logger)

	logger.Info("starting processor", slog.String("address", cfg.HTTPServer.Address()))

	go func() {
		if err := http.ListenAndServe(cfg.HTTPServer.Address(), appRouter); err != nil {
			logger.Error("failed to start processor", "error", err)
			os.Exit(1)
		}
	}()

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
