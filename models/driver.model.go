package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Driver struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	FullName      string         `gorm:"not null" json:"full_name,omitempty"`
	Phone         string         `json:"phone,omitempty"`
	CCCD          string         `json:"cccd,omitempty"`
	IssueDate     time.Time      `json:"issue_date,omitempty"`
	DateOfBirth   time.Time      `json:"date_of_birth,omitempty"`
	Address       string         `json:"address,omitempty"`
	LicenseNumber string         `json:"license_number,omitempty"`
	LicenseExpiry time.Time      `json:"license_expiry,omitempty"`
	ContractorID  uuid.UUID      `gorm:"type:uuid" json:"contractor_id,omitempty"`
	FixedSalary   *float64       `gorm:"not null;default:0" json:"fixed_salary,omitempty"`
	Note          string         `json:"note,omitempty"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Association with Contractor
	Contractor Contractor `gorm:"foreignKey:ContractorID" json:"-"`
}

type CreateDriverRequest struct {
	FullName      string    `json:"full_name" binding:"required"`
	Phone         string    `json:"phone"`
	CCCD          string    `json:"cccd"`
	IssueDate     time.Time `json:"issue_date"` // Ensure this is properly passed as time.Time
	DateOfBirth   time.Time `json:"date_of_birth"`
	Address       string    `json:"address"`
	LicenseNumber string    `json:"license_number"`
	LicenseExpiry time.Time `json:"license_expiry"`
	ContractorID  uuid.UUID `json:"contractor_id,omitempty"` // Optional contractor info (ID)
	FixedSalary   *float64  `json:"fixed_salary,omitempty"`
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
	ContractorID  uuid.UUID `json:"contractor_id,omitempty"`
	FixedSalary   *float64  `json:"fixed_salary,omitempty"`
	Note          string    `json:"note,omitempty"`
}
