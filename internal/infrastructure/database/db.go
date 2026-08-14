package database

import (
	"log/slog"
	"github.com/orandin/slog-gorm"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"media-processing-platform/internal/domain"
)

func InitDB(log *slog.Logger) *gorm.DB {
	dsn := "host=localhost user=user password=password dbname=postgres port=5433 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: slogGorm.New(
				slogGorm.WithHandler(log.Handler()),
		),
	})

	if err != nil {
		log.Error("failed to connect to local database", "error", err)
		os.Exit(1)
	}

	log.Info("successfully connected to local database on port 5433")

	log.Info("running database auto-migrations...")

	if err := db.AutoMigrate(&domain.Task{}); err != nil {
		log.Error("failed to run database auto-migrations", "error", err)
		os.Exit(1)
	}

	log.Info("successfully ran database auto-migrations")

	return db
}