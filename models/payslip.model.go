package models

import (
	"time"

	"github.com/google/uuid"
)

type Payslip struct {
	ID                  uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	ContractorID        uuid.UUID  `gorm:"type:uuid;not null" json:"contractor_id"`
	Contractor          Contractor `gorm:"foreignKey:ContractorID;references:ID" json:"contractor"`
	DriverID            uuid.UUID  `gorm:"type:uuid;not null" json:"driver_id"`
	Driver              Driver     `gorm:"foreignKey:DriverID" json:"driver"`
	TotalTrips          *float64   `gorm:"not null;default:0" json:"total_trips"`
	TakeCareTruckSalary *float64   `gorm:"not null;default:0" json:"take_care_truck_salary"`
	AllowanceSunday     *float64   `gorm:"not null;default:0" json:"allowance_sunday_salary"`
	AllowanceDaily      *float64   `gorm:"not null;default:0" json:"allowance_daily_salary"`
	AllowancePhone      *float64   `gorm:"not null;default:0" json:"allowance_phone_salary"`
	PointSalary         *float64   `gorm:"not null;default:0" json:"point_salary"`
	TripSalary          *float64   `gorm:"not null;default:0" json:"trip_salary"`
	MealFee             *float64   `gorm:"not null;default:0" json:"meal_fee"`
	DailySalary         *float64   `gorm:"not null;default:0" json:"daily_salary"`
	KPISalary           *float64   `gorm:"not null;default:0" json:"kpi_salary"`
	LoadingSalary       *float64   `gorm:"not null;default:0" json:"loading_salary"`
	ParkingFee          *float64   `gorm:"not null;default:0" json:"parking_fee"`
	StandbyFee          *float64   `gorm:"not null;default:0" json:"standby_fee"`
	OtherSalary         *float64   `gorm:"not null;default:0" json:"other_salary"`
	OutsideOilFee       *float64   `gorm:"not null;default:0" json:"outside_oil_fee"`
	OilFee              *float64   `gorm:"not null;default:0" json:"oil_fee"`
	ChargeFee           *float64   `gorm:"not null;default:0" json:"charge_fee"`
	FinalSalary         *float64   `gorm:"not null;default:0" json:"final_salary"`
	DepositSalary       *float64   `gorm:"not null;default:0" json:"deposit_salary"`
	Year                int        `gorm:"not null" json:"year"`
	Month               int        `gorm:"not null" json:"month"`
	Notes               string     `gorm:"type:text" json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
