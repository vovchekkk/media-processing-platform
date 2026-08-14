package repository

import (
	"github.com/google/uuid"

	"media-processing-platform/internal/domain"
)

type Task interface {
	CreateTask(task *domain.Task) error

	GetTaskStatusByID(id uuid.UUID) (domain.TaskStatus, error)

	GetTaskResultByID(id uuid.UUID) (string, error)
}