package postgres

import (
	"context"

	"gorm.io/gorm"

	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
)

type gormUser struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repository.User {
	return &gormUser{db: db}
}

var _ repository.User = (*gormUser)(nil)

func (r *gormUser) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *gormUser) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User

	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}