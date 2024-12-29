package models

import (
	"time"

	"github.com/google/uuid"
)

type Truck struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	LicensePlate string    `gorm:"not null" json:"license_plate,omitempty"`
	Capacity     float64   `gorm:"not null" json:"capacity,omitempty"`
	Length       float64   `gorm:"not null" json:"length,omitempty"`
	Width        float64   `gorm:"not null" json:"width,omitempty"`
	Height       float64   `gorm:"not null" json:"height,omitempty"`
	Volume       float64   `gorm:"not null" json:"volume,omitempty"`
	VehicleType  string    `gorm:"not null" json:"vehicle_type,omitempty"`
	Brand        string    `gorm:"not null" json:"brand,omitempty"`          // New field for Brand
	ContractorID uuid.UUID `gorm:"type:uuid" json:"contractor_id,omitempty"` // Foreign key to Contractor
	Note         string    `json:"note,omitempty"`
	Status       string    `gorm:"not null;default:'active'" json:"status,omitempty"` // New field for Status
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`

	// Association
	Contractor Contractor `gorm:"foreignKey:ContractorID" json:"contractor,omitempty"` // Add Contractor relationship
}

type CreateTruckRequest struct {
	LicensePlate string  `json:"license_plate" binding:"required"`
	Capacity     float64 `json:"capacity" binding:"required"`
	Length       float64 `json:"length" binding:"required"`
	Width        float64 `json:"width" binding:"required"`
	Height       float64 `json:"height" binding:"required"`
	Volume       float64 `json:"volume" binding:"required"`
	VehicleType  string  `json:"vehicle_type" binding:"required"`
	Brand        string  `json:"brand" binding:"required"` // Added field for brand
	Contractor   string  `json:"contractor,omitempty"`     // Optional contractor info (ID)
	Note         string  `json:"note,omitempty"`
	Status       string  `json:"status,omitempty"` // Optional status field
}

type UpdateTruckRequest struct {
	LicensePlate string  `json:"license_plate,omitempty"`
	Capacity     float64 `json:"capacity,omitempty"`
	Length       float64 `json:"length,omitempty"`
	Width        float64 `json:"width,omitempty"`
	Height       float64 `json:"height,omitempty"`
	Volume       float64 `json:"volume,omitempty"`
	VehicleType  string  `json:"vehicle_type,omitempty"`
	Brand        string  `json:"brand,omitempty"`      // Updated brand field
	Contractor   string  `json:"contractor,omitempty"` // Optional contractor info (ID)
	Note         string  `json:"note,omitempty"`
	Status       string  `json:"status,omitempty"` // Optional status field
}
