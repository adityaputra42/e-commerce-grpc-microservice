package model

import "time"

type UserAddresses struct {
	ID             string    `gorm:"primaryKey;column:id"`
	Username       string    `gorm:"primaryKey;column:username"`
	Label          string    `gorm:"column:label"`
	RecipientName  string    `gorm:"column:recipient_name"`
	RecipientPhone string    `gorm:"column:recipient_phone"`
	AddressLine    string    `gorm:"column:address_line"`
	City           string    `gorm:"column:city"`
	Province       string    `gorm:"column:province"`
	PostalCode     string    `gorm:"column:postal_code"`
	IsSelected     bool      `gorm:"column:is_selected"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *UserAddresses) TableName() string {
	return "users"
}
