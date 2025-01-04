package models

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	FullName      string    `gorm:"not null" json:"full_name,omitempty"`
	Phone         string    `gorm:"not null" json:"phone,omitempty"`
	CCCD          string    `gorm:"not null" json:"cccd,omitempty"`
	IssueDate     time.Time `gorm:"not null" json:"issue_date,omitempty"`
	DateOfBirth   time.Time `gorm:"not null" json:"date_of_birth,omitempty"`
	Address       string    `gorm:"not null" json:"address,omitempty"`
	LicenseNumber string    `gorm:"not null" json:"license_number,omitempty"`
	LicenseExpiry time.Time `gorm:"not null" json:"license_expiry,omitempty"`
	ContractorID  uuid.UUID `gorm:"type:uuid" json:"contractor_id,omitempty"` // foreign key to Contractor
	Note          string    `json:"note,omitempty"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at,omitempty"`

	// Association
	Contractor Contractor `gorm:"foreignKey:ContractorID" json:"contractor,omitempty"` // Add Contractor relationship
}

type CreateDriverRequest struct {
	FullName      string    `json:"full_name" binding:"required"`
	Phone         string    `json:"phone" binding:"required"`
	CCCD          string    `json:"cccd" binding:"required"`
	IssueDate     time.Time `json:"issue_date" binding:"required"`
	DateOfBirth   time.Time `json:"date_of_birth" binding:"required"`
	Address       string    `json:"address" binding:"required"`
	LicenseNumber string    `json:"license_number" binding:"required"`
	LicenseExpiry time.Time `json:"license_expiry" binding:"required"`
	Contractor    string    `json:"contractor,omitempty"` // optional contractor info (ID)
	Note          string    `json:"note,omitempty"`
}

type UpdateDriverRequest struct {
	FullName      string    `json:"full_name,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	CCCD          string    `json:"cccd,omitempty"`
	IssueDate     time.Time `json:"issue_date,omitempty"`
	DateOfBirth   time.Time `json:"date_of_birth,omitempty"`
	Address       string    `json:"address,omitempty"`
	LicenseNumber string    `json:"license_number,omitempty"`
	LicenseExpiry time.Time `json:"license_expiry,omitempty"`
	Contractor    string    `json:"contractor,omitempty"` // optional contractor info (ID)
	Note          string    `json:"note,omitempty"`
}