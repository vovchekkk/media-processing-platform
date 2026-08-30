package domain

import (
	"image"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID     uuid.UUID    `gorm:"type:uuid;primaryKey;"`
	Filter *ImageFilter `gorm:"serializer:json"`
	Image  image.Image  `gorm:"type:text"`
}

type ImageFilter struct {
	Name       string         `json:"name"`
	Parameters map[string]any `json:"parameters"`
}
