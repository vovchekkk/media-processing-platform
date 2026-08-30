package repository

import (
	"context"

	"media-processing-platform/server/internal/domain"
)

type User interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}