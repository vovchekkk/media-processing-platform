package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	StatusInProgress    TaskStatus = "in_progress"
	StatusReady TaskStatus = "ready"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"` 
	Status    string    `gorm:"type:varchar(20);default:'in_progress'"`
	Result    string    `gorm:"type:text"`
}