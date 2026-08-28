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

func (r *gormTask) GetTaskStatusByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.TaskStatus, error) {
	var task domain.Task

	err := r.db.WithContext(ctx).Select("status").First(&task, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return "", err
	}

	return domain.TaskStatus(task.Status), nil
}

func (r *gormTask) GetTaskResultByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error) {
	var task domain.Task

	err := r.db.WithContext(ctx).Select("result").First(&task, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return "", err
	}

	return task.Result, nil
}

func (r *gormTask) UpdateTaskStatusByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.TaskStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Task{}).Where("id = ? AND user_id = ?", id, userID).Update("status", status).Error
}

func (r *gormTask) SetTaskResultByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID, result string) error {
	return r.db.WithContext(ctx).Model(&domain.Task{}).Where("id = ? AND user_id = ?", id, userID).Update("result", result).Error
}
