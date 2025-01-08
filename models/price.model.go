package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CreatePricing is used to create a new pricing record
type CreatePricing struct {
	ContractorID uuid.UUID     `gorm:"type:uuid" json:"contractor_id,omitempty"`
	FileName     string        `gorm:"not null" json:"file_name"`
	Prices       []PriceDetail `gorm:"-" json:"prices"` // This is the list of price details provided in the request
}

// Pricing represents the pricing structure in the system
type Pricing struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ContractorID uuid.UUID     `gorm:"type:uuid" json:"contractor_id,omitempty"`
	FileName     string        `gorm:"not null" json:"file_name"`
	PriceDetails []PriceDetail `gorm:"foreignkey:PricingID" json:"price_details"` // One-to-many relationship
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type PriceDetail struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	PricingID    uuid.UUID `gorm:"type:uuid" json:"pricing_id"` // Foreign key to Pricing
	FromCity     string    `json:"from_city"`
	FromDistrict string    `json:"from_district"`
	ToCity       string    `json:"to_city"`
	ToDistrict   string    `json:"to_district"`
	Notes        string    `json:"notes"`
	WeightPrices JSONBMap  `gorm:"type:jsonb" json:"weight_prices"` // Use custom type to handle map
}

// Custom type to handle map in JSONB format
type JSONBMap map[string]float64

// Implement GORM's Scanner and Valuer interfaces for the JSONBMap custom type

// Scan implements the Scanner interface to handle the JSONB type
func (j *JSONBMap) Scan(value interface{}) error {
	// If value is nil, initialize as an empty map
	if value == nil {
		*j = JSONBMap{}
		return nil
	}
	// Otherwise, unmarshal the value to JSONBMap
	return json.Unmarshal(value.([]byte), j)
}

// Value implements the Valuer interface to store the map as JSONB in PostgreSQL
func (j JSONBMap) Value() (driver.Value, error) {
	// Marshal the map to JSON and store it as []byte
	return json.Marshal(j)
}
