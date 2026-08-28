package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Username string    `gorm:"type:varchar(255);unique;not null"`
	Password string    `gorm:"type:varchar(255);not null"`
}
