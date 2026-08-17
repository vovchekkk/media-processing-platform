package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
)

type gormSession struct {
	db *gorm.DB
}

func NewGormSessionRepository(db *gorm.DB) repository.Session {
	return &gormSession{db: db}
}

var _ repository.Session = (*gormSession)(nil)

func (r *gormSession) CreateSession(ctx context.Context, session *domain.Session) error {
	panic("unimplemented")
}

func (r *gormSession) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	panic("unimplemented")
}

func (r *gormSession) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	panic("unimplemented")
}