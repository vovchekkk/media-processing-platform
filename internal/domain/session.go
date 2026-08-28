package domain

import "github.com/google/uuid"

type Session struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	User User `gorm:"foreignKey:UserID;references:ID;constraints:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}
