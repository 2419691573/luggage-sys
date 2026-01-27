package models

import (
	"time"
)

// StoredLog 寄存记录
type StoredLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HotelID   uint      `gorm:"not null" json:"hotel_id"`
	LuggageID uint      `gorm:"not null" json:"luggage_id"`
	GuestName string    `gorm:"not null" json:"guest_name"`
	Status    string    `gorm:"not null" json:"status"`
	StoredAt  time.Time `gorm:"autoCreateTime" json:"stored_at"`
}

func (StoredLog) TableName() string {
	return "stored_logs"
}

// UpdatedLog 修改记录
type UpdatedLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HotelID   uint      `gorm:"not null" json:"hotel_id"`
	LuggageID uint      `gorm:"not null" json:"luggage_id"`
	UpdatedBy string    `gorm:"not null" json:"updated_by"`
	OldData   string    `gorm:"type:text" json:"old_data"`
	NewData   string    `gorm:"type:text" json:"new_data"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

func (UpdatedLog) TableName() string {
	return "updated_logs"
}

// RetrievedLog 取出记录
type RetrievedLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HotelID     uint      `gorm:"not null" json:"hotel_id"`
	LuggageID   uint      `gorm:"not null" json:"luggage_id"`
	GuestName   string    `gorm:"not null" json:"guest_name"`
	RetrievedBy string    `gorm:"not null" json:"retrieved_by"`
	RetrievedAt time.Time `gorm:"autoCreateTime" json:"retrieved_at"`
}

func (RetrievedLog) TableName() string {
	return "retrieved_logs"
}
