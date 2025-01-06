package models

import (
	"time"

	"github.com/google/uuid"
)

type Pricing struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`                                                  // ID của bản ghi
	ContractorID         uuid.UUID      `gorm:"type:uuid" json:"contractor_id,omitempty"`                              // Foreign key for Contractor                                  // ID của nhà thầu
	Contractor           Contractor     `gorm:"foreignKey:ContractorID;constraint:OnDelete:CASCADE" json:"contractor"` // Mối quan hệ với nhà thầu
	FromLocationCity     string         `gorm:"not null" json:"from_location_city"`                                    // Thành phố lấy hàng
	FromLocationDistrict string         `gorm:"not null" json:"from_location_district"`                                // Quận/Huyện lấy hàng
	ToLocationCity       string         `gorm:"not null" json:"to_location_city"`                                      // Thành phố trả hàng
	ToLocationDistrict   string         `gorm:"not null" json:"to_location_district"`                                  // Quận/Huyện trả hàng
	Prices               map[string]int `gorm:"type:jsonb;not null" json:"prices"`                                     // Giá theo các loại tải trọng (theo tấn, khối)
	Note                 string         `json:"note"`                                                                  // Ghi chú về bảng giá
	Version              string         `gorm:"not null" json:"version"`                                               // Phiên bản của bảng giá (dạng timestamp)
	IsCurrent            bool           `gorm:"default:true" json:"is_current"`                                        // Đánh dấu bảng giá hiện tại
	CreatedAt            time.Time      `json:"created_at"`                                                            // Thời gian tạo
	UpdatedAt            time.Time      `json:"updated_at"`                                                            // Thời gian cập nhật
}
