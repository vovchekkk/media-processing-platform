package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"media-processing-platform/server/internal/domain"
)

type Session interface {
	CreateSession(ctx context.Context, session *domain.Session, ttl time.Duration) error
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
}
