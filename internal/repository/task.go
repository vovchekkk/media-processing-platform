package repository

import (
	"media-processing-platform/internal/domain"
)

type Task interface {
	CreateTask(task *domain.Task) error

	GetTaskStatusByID(id string) (string, error)

	GetTaskResultByID(id string) (string, error)
}