package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ContractorID     uuid.UUID  `gorm:"type:uuid;not null" json:"contractor_id"`
	Contractor       Contractor `gorm:"foreignKey:ContractorID;references:ID" json:"contractor"`
	OrderTime        time.Time  `json:"order_time"`
	CompanyName      string     `gorm:"size:255;not null" json:"company_name"`
	DriverID         uuid.UUID  `gorm:"type:uuid;not null" json:"driver_id"`
	Driver           Driver     `gorm:"foreignKey:DriverID" json:"driver"`
	TruckID          uuid.UUID  `gorm:"type:uuid;not null" json:"truck_id"`
	Truck            Truck      `gorm:"foreignKey:TruckID" json:"truck"`
	PriceID          uuid.UUID  `gorm:"type:uuid;not null" json:"price_id"`
	Pricing          Pricing    `gorm:"foreignKey:PriceID" json:"pricing"`
	PickupProvince   string     `gorm:"size:50;not null" json:"pickup_province"`
	PickupDistrict   string     `gorm:"size:50" json:"pickup_district"`
	DeliveryProvince string     `gorm:"size:50;not null" json:"delivery_province"`
	DeliveryDistrict string     `gorm:"size:50" json:"delivery_district"`
	Unit             string     `gorm:"size:20;not null" json:"unit"`
	Weight           *float64   `json:"weight"`
	Volume           *float64   `json:"volume"`
	TripCount        int        `gorm:"not null" json:"trip_count"`
	FreightCharge    *float64   `gorm:"not null" json:"freight_charge"`
	PointFee         *float64   `json:"point_fee"`
	PointCount       *int       `json:"point_count"`
	RefundFee        *float64   `json:"refund_fee"`
	LoadingFee       *float64   `json:"loading_fee"`
	MealFee          *float64   `json:"meal_fee"`
	StandbyFee       *float64   `json:"standby_fee"`
	ParkingFee       *float64   `json:"parking_fee"`
	OtherFees        []OtherFee `gorm:"foreignKey:OrderID" json:"other_fees"`
	Note             string     `gorm:"type:text" json:"note"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type OtherFee struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	OrderID uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	Name    string    `gorm:"size:255;not null" json:"name"`
	Value   float64   `gorm:"not null" json:"value"`
}
