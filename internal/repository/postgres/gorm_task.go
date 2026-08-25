package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"media-processing-platform/internal/domain"
	"media-processing-platform/internal/repository"
)

type gormTask struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) repository.Task {
	return &gormTask{db: db}
}

var _ repository.Task = (*gormTask)(nil)

func (r *gormTask) CreateTask(ctx context.Context, task *domain.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *gormTask) GetTaskStatusByID(ctx context.Context, id uuid.UUID) (domain.TaskStatus, error) {
	var task domain.Task

	err := r.db.WithContext(ctx).Select("status").First(&task, "id = ?", id).Error
	if err != nil {
		return "", err
	}

	return domain.TaskStatus(task.Status), nil
}

func (r *gormTask) GetTaskResultByID(ctx context.Context, id uuid.UUID) (string, error) {
	var task domain.Task

	err := r.db.WithContext(ctx).Select("result").First(&task, "id = ?", id).Error
	if err != nil {
		return "", err
	}

	return task.Result, nil
}

func (r *gormTask) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *gormTask) SetTaskResult(ctx context.Context, id uuid.UUID, result string) error {
	return r.db.WithContext(ctx).Model(&domain.Task{}).Where("id = ?", id).Update("result", result).Error
}
