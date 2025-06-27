package model

import "time"

type Cars struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Name         string    `gorm:"column:name"`
	Brand        string    `gorm:"column:brand"`
	Model        string    `gorm:"column:model"`
	Year         int32     `gorm:"column:year"`
	Mileage      int32     `gorm:"column:mileage"`
	Transmission string    `gorm:"column:transmission"`
	FuelType     string    `gorm:"column:fuel_type"`
	Location     string    `gorm:"column:location"`
	Description  string    `gorm:"column:description"`
	Images       []string  `gorm:"column:images"`
	Price        float64   `gorm:"column:price"`
	Currency     string    `gorm:"column:currency"`
	IsSold       bool      `gorm:"column:is_sold"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *Cars) TableName() string {
	return "cars"
}
