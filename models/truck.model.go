package models

import (
	"time"

	"github.com/google/uuid"
)

// Truck represents a truck entity with associated details.
type Truck struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	LicensePlate string    `gorm:"not null" json:"license_plate,omitempty"`
	Capacity     float64   `gorm:"not null" json:"capacity,omitempty"`
	Length       float64   `gorm:"not null" json:"length,omitempty"`
	Width        float64   `gorm:"not null" json:"width,omitempty"`
	Height       float64   `gorm:"not null" json:"height,omitempty"`
	Volume       float64   `gorm:"not null" json:"volume,omitempty"`
	Brand        string    `json:"brand,omitempty"`                          // Truck brand
	ContractorID uuid.UUID `gorm:"type:uuid" json:"contractor_id,omitempty"` // Foreign key for Contractor
	Note         string    `json:"note,omitempty"`
	Status       string    `gorm:"not null;default:'active'" json:"status,omitempty"` // Default to 'active'
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`

	// Association
	Contractor Contractor `gorm:"foreignKey:ContractorID" json:"contractor,omitempty"` // Contractor relationship
}

// CreateTruckRequest represents the request payload for creating a truck.
type CreateTruckRequest struct {
	LicensePlate string    `json:"license_plate" binding:"required"`
	Capacity     float64   `json:"capacity" binding:"required"`
	Length       float64   `json:"length" binding:"required"`
	Width        float64   `json:"width" binding:"required"`
	Height       float64   `json:"height" binding:"required"`
	Volume       float64   `json:"volume" binding:"required"`
	Brand        string    `json:"brand,omitempty"`         // Truck brand
	ContractorID uuid.UUID `json:"contractor_id,omitempty"` // Optional contractor ID
	Note         string    `json:"note,omitempty"`
	Status       string    `json:"status,omitempty"` // Optional status field
}

// UpdateTruckRequest represents the request payload for updating truck details.
type UpdateTruckRequest struct {
	LicensePlate string    `json:"license_plate,omitempty"`
	Capacity     float64   `json:"capacity,omitempty"`
	Length       float64   `json:"length,omitempty"`
	Width        float64   `json:"width,omitempty"`
	Height       float64   `json:"height,omitempty"`
	Volume       float64   `json:"volume,omitempty"`
	Brand        string    `json:"brand,omitempty"`         // Truck brand
	ContractorID uuid.UUID `json:"contractor_id,omitempty"` // Optional contractor ID
	Note         string    `json:"note,omitempty"`
	Status       string    `json:"status,omitempty"` // Optional status field
}
