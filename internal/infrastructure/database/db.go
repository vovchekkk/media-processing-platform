package database

import (
	"log/slog"
	"github.com/orandin/slog-gorm"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(logger *slog.Logger) *gorm.DB {
	dsn := "host=localhost user=user password=password dbname=postgres port=5433 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: slogGorm.New(
				slogGorm.WithHandler(logger.Handler()),
		),
	})

	if err != nil {
		logger.Error("failed to connect to local database", "error", err)
		os.Exit(1)
	}

	logger.Info("successfully connected to local database on port 5433")
	return db
}