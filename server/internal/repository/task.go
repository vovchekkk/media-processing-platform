package repository

import (
	"context"

	"github.com/google/uuid"

	"media-processing-platform/server/internal/domain"
)

type Task interface {
	CreateTask(ctx context.Context, task *domain.Task) error
	GetTaskStatusByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (domain.TaskStatus, error)
	GetTaskResultByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (string, error)
	UpdateTaskStatusByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.TaskStatus) error
	SetTaskResultByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID, result string) error
}
