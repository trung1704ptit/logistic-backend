package models

import (
	"time"

	"github.com/google/uuid"
)

type Setting struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Settings  JSONBMap  `gorm:"type:jsonb" json:"settings"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}
