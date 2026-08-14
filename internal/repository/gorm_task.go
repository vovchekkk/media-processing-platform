package repository

import (
	"gorm.io/gorm"
	"github.com/google/uuid"

	"media-processing-platform/internal/domain"
)

type gormTask struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) Task {
	return &gormTask{db: db}
}

var _ Task = (*gormTask)(nil)

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