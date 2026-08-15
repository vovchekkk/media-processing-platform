package task

import (
	"media-processing-platform/internal/domain"

	"github.com/google/uuid"
)

type CreateResponse struct {
	ID     uuid.UUID `json:"task_id"`
}

type GetStatusResponse struct {
	TaskStatus domain.TaskStatus `json:"status"`
}

type GetResultResponse struct {
	Result string    `json:"result"`
}