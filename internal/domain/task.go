package domain

import "github.com/google/uuid"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;"` 
	Status    string    `gorm:"type:varchar(20);default:'pending'"`
	Result    string    `gorm:"type:text"`
}