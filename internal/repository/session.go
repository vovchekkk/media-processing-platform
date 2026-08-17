package repository

import (
	"context"

	"github.com/google/uuid"

	"media-processing-platform/internal/domain"
)

type Session interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
}