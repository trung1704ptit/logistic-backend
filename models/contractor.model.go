package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contractor struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string         `gorm:"not null" json:"name,omitempty"`
	Phone     string         `gorm:"not null" json:"phone,omitempty"`
	Address   string         `gorm:"not null" json:"address,omitempty"`
	Note      string         `json:"note,omitempty"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Type      string         `gorm:"not null;default:external" json:"type,omitempty"`

	Drivers  []Driver  `gorm:"foreignKey:ContractorID" json:"drivers,omitempty"` // One-to-many relation with Driver
	Trucks   []Truck   `gorm:"foreignKey:ContractorID" json:"trucks,omitempty"`  // One-to-many relation with Truck
	Pricings []Pricing `gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE;" json:"pricings,omitempty"`
}

type CreateContractorRequest struct {
	Name      string    `json:"name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Note      string    `json:"note,omitempty"`
	Type      string    `gorm:"not null;default:external" json:"type,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdateContractor struct {
	Name      string    `json:"name,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	Note      string    `json:"note,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
