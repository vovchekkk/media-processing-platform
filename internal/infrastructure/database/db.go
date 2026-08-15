package database

import (
	"log/slog"
	"os"

	"github.com/orandin/slog-gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-processing-platform/internal/config"
	"media-processing-platform/internal/domain"
)

func InitDB(dbConfig config.DatabaseConfig, log *slog.Logger) *gorm.DB {
	dsn := dbConfig.DSN()

	log.Info(
		"database config",
    	"host", dbConfig.Host,
    	"port", dbConfig.Port,
    	"dbname", dbConfig.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: slogGorm.New(
				slogGorm.WithHandler(log.Handler()),
		),
	})

	if err != nil {
		log.Error("failed to connect to database", "host", dbConfig.Host, "port", dbConfig.Port, "error", err)
		os.Exit(1)
	}

	log.Info("successfully connected to database", "host", dbConfig.Host, "port", dbConfig.Port)

	log.Info("running database auto-migrations...")

	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		log.Error("failed to run database auto-migrations", "error", err)
		os.Exit(1)
	}

	log.Info("successfully ran database auto-migrations")

	return db
}