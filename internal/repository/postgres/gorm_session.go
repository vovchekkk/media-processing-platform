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
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *gormSession) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	var session domain.Session

	err := r.db.WithContext(ctx).First(&session, "id = ?", sessionID).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *gormSession) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Session{}, "id = ?", sessionID).Error
}
