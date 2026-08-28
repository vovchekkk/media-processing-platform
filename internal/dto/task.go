package dto

import (
	"media-processing-platform/internal/domain"
)

type Task struct {
	Filter domain.ImageFilter `json:"filter"`
	Image  string             `json:"image"`
}
