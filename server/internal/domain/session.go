package domain

import "github.com/google/uuid"

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
