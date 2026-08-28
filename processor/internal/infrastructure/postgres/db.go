package database

import (
	"log/slog"
	"os"

	slogGorm "github.com/orandin/slog-gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-processing-platform/processor/internal/config"
)

func InitDB(cfg config.DatabaseConfig, log *slog.Logger) *gorm.DB {
	dsn := cfg.DSN()

	log.Info(
		"database config",
		"host", cfg.Host,
		"port", cfg.Port,
		"dbname", cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: slogGorm.New(
			slogGorm.WithHandler(log.Handler()),
		),
	})

	if err != nil {
		log.Error("failed to connect to database", "host", cfg.Host, "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	log.Info("successfully connected to database", "host", cfg.Host, "port", cfg.Port)

	return db
}
