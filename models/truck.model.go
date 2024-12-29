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
	Contractor   uuid.UUID `gorm:"type:uuid" json:"contractor,omitempty"`
	Note         string    `json:"note,omitempty"`
	CreatedAt    time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt    time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateTruckRequest struct {
	LicensePlate string  `json:"license_plate" binding:"required"`
	Capacity     float64 `json:"capacity" binding:"required"`
	Length       float64 `json:"length" binding:"required"`
	Width        float64 `json:"width" binding:"required"`
	Height       float64 `json:"height" binding:"required"`
	Volume       float64 `json:"volume" binding:"required"`
	VehicleType  string  `json:"vehicle_type" binding:"required"`
	Contractor   string  `json:"contractor,omitempty"`
	Note         string  `json:"note,omitempty"`
}

type UpdateTruckRequest struct {
	LicensePlate string  `json:"license_plate,omitempty"`
	Capacity     float64 `json:"capacity,omitempty"`
	Length       float64 `json:"length,omitempty"`
	Width        float64 `json:"width,omitempty"`
	Height       float64 `json:"height,omitempty"`
	Volume       float64 `json:"volume,omitempty"`
	VehicleType  string  `json:"vehicle_type,omitempty"`
	Contractor   string  `json:"contractor,omitempty"`
	Note         string  `json:"note,omitempty"`
}
