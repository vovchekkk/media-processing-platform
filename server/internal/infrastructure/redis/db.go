package redis

import (
	"context"
	"log/slog"
	"media-processing-platform/server/internal/config"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitDB(ctx context.Context, cfg config.RedisConfig, log *slog.Logger) *redis.Client {
	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Address(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.Connection.MaxRetries,
		DialTimeout:  cfg.Connection.DialTimeout,
		ReadTimeout:  cfg.Connection.Timeout,
		WriteTimeout: cfg.Connection.Timeout,
	})

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Connection.Timeout)
	defer cancel()

	if err := db.Ping(pingCtx).Err(); err != nil {
		log.Error("failed to connect to redis", "error", err)
		os.Exit(1)
	}

	log.Info("successfully connected to redis", "host", cfg.Host, "port", cfg.Port)

	return db
}
