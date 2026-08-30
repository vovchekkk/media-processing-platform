package repository

import (
	"context"

	"github.com/google/uuid"

	"media-processing-platform/processor/internal/domain"
)

type Task interface {
	UpdateTaskStatus(ctx context.Context, id uuid.UUID, status domain.TaskStatus) error
	UpdateTaskResult(ctx context.Context, id uuid.UUID, result string) error
	UpdateTaskStatusAndResult(ctx context.Context, id uuid.UUID, status domain.TaskStatus, result string) error
}
