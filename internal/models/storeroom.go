package models

import (
	"time"

	"gorm.io/gorm"
)

type Storeroom struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	HotelID           uint           `gorm:"not null" json:"hotel_id"`
	Name              string         `gorm:"not null" json:"name"`
	Location          string         `json:"location"`
	Capacity          int            `gorm:"not null" json:"capacity"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	StoredCount       int            `gorm:"-" json:"stored_count"`      // 计算字段，不存储
	RemainingCapacity int            `gorm:"-" json:"remaining_capacity"` // 计算字段，不存储
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Storeroom) TableName() string {
	return "storerooms"
}
