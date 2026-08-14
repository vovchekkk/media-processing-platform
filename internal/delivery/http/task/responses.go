package task

import (
	"github.com/google/uuid"
	resp "media-processing-platform/internal/delivery/http/shared"
)

type CreateResponse struct {
	ID     uuid.UUID `json:"id"`
	resp.Response
}