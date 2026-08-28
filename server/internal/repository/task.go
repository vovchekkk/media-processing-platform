package repository

import (
	"context"

	"github.com/google/uuid"

	"media-processing-platform/server/internal/domain"
)

type Task interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.TaskStatus, error)
	GetTaskResult(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error)
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.TaskStatus) error
	UpdateTaskResult(ctx context.Context, id uuid.UUID, userID uuid.UUID, result string) error
	UpdateTaskStatusAndResult(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.TaskStatus, result string) error
}
