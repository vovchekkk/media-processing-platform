package repository

import (
	"context"

	"github.com/google/uuid"

	"media-processing-platform/internal/domain"
)

type Task interface {
	CreateTask(ctx context.Context, task *domain.Task) error

	GetTaskStatusByID(ctx context.Context, id uuid.UUID) (domain.TaskStatus, error)

	GetTaskResultByID(ctx context.Context, id uuid.UUID) (string, error)

	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error

	SetTaskResult(ctx context.Context, id uuid.UUID, result string) error
}
