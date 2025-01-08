package models

import (
	"time"

	"github.com/google/uuid"
)

type PriceDetail struct {
	ID           uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	FromProvince string             `json:"from_province"`
	FromDistrict string             `json:"from_district"`
	ToProvince   string             `json:"to_province"`
	ToDistrict   string             `json:"to_district"`
	Notes        string             `json:"notes"`
	WeightPrices map[string]float64 `json:"weight_prices"`
}

type Pricing struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ContractorID uuid.UUID     `gorm:"type:uuid" json:"contractor_id,omitempty"`
	FileName     string        `gorm:"not null" json:"file_name"`
	Prices       []PriceDetail `gorm:"type:jsonb;not null" json:"prices"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
