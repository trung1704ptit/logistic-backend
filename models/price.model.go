package models

import (
	"time"

	"github.com/google/uuid"
)

type Pricing struct {
	ID           uint                   `gorm:"primaryKey" json:"id"`                                                  // ID của bản ghi
	ContractorID uuid.UUID              `gorm:"type:uuid" json:"contractor_id,omitempty"`                              // Foreign key for Contractor                                  // ID của nhà thầu
	Contractor   Contractor             `gorm:"foreignKey:ContractorID;constraint:OnDelete:CASCADE" json:"contractor"` // Mối quan hệ với nhà thầu
	FromCity     string                 `gorm:"not null" json:"from_city"`                                             // Thành phố lấy hàng
	FromDistrict string                 `gorm:"not null" json:"from_district"`                                         // Quận/Huyện lấy hàng
	ToCity       string                 `gorm:"not null" json:"to_city"`                                               // Thành phố trả hàng
	ToDistrict   string                 `gorm:"not null" json:"to_district"`                                           // Quận/Huyện trả hàng
	Prices       map[string]interface{} `gorm:"type:jsonb;not null" json:"prices"`                                     // Giá theo các loại tải trọng (theo tấn, khối)
	Note         string                 `json:"note"`                                                                  // Ghi chú về bảng giá
	FileName     string                 `gorm:"not null" json:"file_name"`                                             // Phiên bản của bảng giá (dạng timestamp)
	CreatedAt    time.Time              `json:"created_at"`                                                            // Thời gian tạo
	UpdatedAt    time.Time              `json:"updated_at"`                                                            // Thời gian cập nhật
}
