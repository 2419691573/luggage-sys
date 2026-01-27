package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Role     string `gorm:"type:varchar(32);not null;default:staff" json:"role"`
	HotelID  uint   `gorm:"not null" json:"hotel_id"`
}

func (User) TableName() string {
	return "users"
}
