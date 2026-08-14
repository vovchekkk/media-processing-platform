package task

import (
	resp "media-processing-platform/internal/delivery/http/shared"
	"media-processing-platform/internal/domain"

	"github.com/google/uuid"
)

type CreateResponse struct {
	ID     uuid.UUID `json:"task_id"`
	resp.Response
}

type GetStatusResponse struct {
	TaskStatus domain.TaskStatus `json:"task_status"`
	resp.Response
}

type GetResultResponse struct {
	Result string    `json:"result"`
	resp.Response
}