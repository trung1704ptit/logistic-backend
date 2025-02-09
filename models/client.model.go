package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Client struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string         `gorm:"not null" json:"name,omitempty"`
	Phone     string         `gorm:"not null" json:"phone,omitempty"`
	Address   string         `gorm:"not null" json:"address,omitempty"`
	Note      string         `json:"note,omitempty"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Pricings []Pricing `gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE;" json:"pricings,omitempty"`
}

type CreateClientRequest struct {
	Name      string    `json:"name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdateClient struct {
	Name      string    `json:"name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
