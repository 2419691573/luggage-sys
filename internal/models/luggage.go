package models

import (
	"time"

	"gorm.io/gorm"
)

type Luggage struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	GuestName     string    `gorm:"not null" json:"guest_name"`
	StaffName     string    `gorm:"not null" json:"staff_name"`
	ContactPhone  string    `json:"contact_phone"`
	ContactEmail  string    `json:"contact_email"`
	Description   string    `json:"description"`
	Quantity      int       `gorm:"default:1" json:"quantity"`
	SpecialNotes  string    `json:"special_notes"`
	// PhotoURLs stores multiple image URLs (recommended).
	// Note: keep PhotoURL for backward compatibility (first image).
	PhotoURLs     StringSlice `gorm:"type:json" json:"photo_urls"`
	PhotoURL      string    `json:"photo_url"`
	StoreroomID   uint      `gorm:"not null" json:"storeroom_id"`
	Storeroom     Storeroom `gorm:"foreignKey:StoreroomID" json:"-"`
	RetrievalCode string    `gorm:"type:varchar(32);index;not null" json:"retrieval_code"` // 普通索引，允许多个行李共用同一个取件码
	Status        string    `gorm:"not null;default:stored" json:"status"` // stored, retrieved
	StoredAt      time.Time `gorm:"autoCreateTime" json:"stored_at"`
	RetrievedAt   *time.Time `json:"retrieved_at,omitempty"`
	RetrievedBy   string    `json:"retrieved_by,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Luggage) TableName() string {
	return "luggages"
}
