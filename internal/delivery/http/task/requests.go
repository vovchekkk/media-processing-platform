package task

import "github.com/google/uuid"

type CreateRequest struct{}

type GetStatusRequest struct {
	ID uuid.UUID `json:"task_id"`
}

type GetResultRequest struct {
	ID uuid.UUID `json:"task_id"`
}
