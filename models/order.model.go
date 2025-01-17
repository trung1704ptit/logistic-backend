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
	PickupProvince   string     `gorm:"size:50;not null" json:"pickup_province"`
	PickupDistrict   string     `gorm:"size:50" json:"pickup_district"`
	DeliveryProvince string     `gorm:"size:50;not null" json:"delivery_province"`
	DeliveryDistrict string     `gorm:"size:50" json:"delivery_district"`
	Unit             string     `gorm:"size:20;not null" json:"unit"`
	PackageWeight    *float64   `json:"package_weight"`
	PackageVolume    *float64   `json:"package_volumn"`
	TripCount        int        `gorm:"not null;default:1" json:"trip_count"`
	TripSalary       *float64   `gorm:"not null;default:0" json:"trip_salary"`
	DailySalary      *float64   `gorm:"not null;default:0" json:"daily_salary"`
	PointCount       *int       `gorm:"not null;default:1" json:"point_count"`
	PointSalary      *float64   `gorm:"not null;default:0" json:"point_salary"`
	RefundFee        *float64   `gorm:"not null;default:0" json:"recovery_fee"`
	LoadingSalary    *float64   `gorm:"not null;default:0" json:"loading_salary"`
	MealFee          *float64   `gorm:"not null;default:0" json:"meal_fee"`
	StandbyFee       *float64   `gorm:"not null;default:0" json:"standby_fee"`
	ParkingFee       *float64   `gorm:"not null;default:0" json:"parking_fee"`
	OtherSalary      *float64   `gorm:"not null;default:0" json:"other_salary"`
	OutsiteOilFee    *float64   `gorm:"not null;default:0" json:"outside_oil_fee"`
	OilFee           *float64   `gorm:"not null;default:0" json:"oil_fee"`
	ChargeFee        *float64   `gorm:"not null;default:0" json:"charge_fee"`
	TotalSalary      *float64   `gorm:"not null;default:0" json:"total_salary"`
	Notes            string     `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}
