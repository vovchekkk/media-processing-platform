package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	StatusInProgress TaskStatus = "in_progress"
	StatusReady      TaskStatus = "ready"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID     uuid.UUID   `gorm:"type:uuid;primaryKey;"`
	UserID uuid.UUID   `gorm:"type:uuid;not null"`
	Filter ImageFilter `gorm:"serializer:json"`
	Image  string      `gorm:"type:text"`
	Status string      `gorm:"type:varchar(255);not null;default:'in_progress'"`
	Result string      `gorm:"type:text"`

	User User `gorm:"foreignKey:UserID;references:ID;constraints:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type ImageFilter struct {
	Name       string         `json:"name"`
	Parameters map[string]any `json:"parameters"`
}
