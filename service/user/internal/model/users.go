package model

import "time"

type Users struct {
	Username    string    `gorm:"primaryKey;column:username"`
	FullName    string    `gorm:"column:full_name"`
	PhoneNumber string    `gorm:"column:phone_number"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *Users) TableName() string {
	return "users"
}
