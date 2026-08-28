package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-processing-platform/processor/internal/domain"
	"media-processing-platform/processor/internal/repository"
)

type gormTask struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) repository.Task {
	return &gormTask{db: db}
}

var _ repository.Task = (*gormTask)(nil)

func (r *gormTask) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	return r.db.WithContext(ctx).
		Table("tasks").
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *gormTask) UpdateTaskResult(ctx context.Context, id uuid.UUID, result string) error {
	return r.db.WithContext(ctx).
		Table("tasks").
		Where("id = ?", id).
		Update("result", result).Error
}

func (r *gormTask) UpdateTaskStatusAndResult(ctx context.Context, id uuid.UUID, status domain.TaskStatus, result string) error {
	return r.db.WithContext(ctx).
		Table("tasks").
		Where("id = ?", id).
		Updates(map[string]any{
			"status": status,
			"result": result,
		}).Error
}
