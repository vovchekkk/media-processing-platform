package postgres

import (
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

func (r *gormTask) CreateTask(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *gormTask) GetTaskStatusByID(id uuid.UUID) (domain.TaskStatus, error) {
	var task domain.Task

	err := r.db.Select("status").First(&task, "id = ?", id).Error
	if err != nil {
		return "", err
	}

	return domain.TaskStatus(task.Status), nil
}

func (r *gormTask) GetTaskResultByID(id uuid.UUID) (string, error) {
	var task domain.Task

	err := r.db.Select("result").First(&task, "id = ?", id).Error
	if err != nil {
		return "", err
	}

	return task.Result, nil
}

func (r *gormTask) UpdateTaskStatus(id uuid.UUID, status domain.TaskStatus) error {
	return r.db.Model(&domain.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *gormTask) SetTaskResult(id uuid.UUID, result string) error {
	return r.db.Model(&domain.Task{}).Where("id = ?", id).Update("result", result).Error
}
