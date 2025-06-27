package model

import "time"

type Cars struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Name         string    `gorm:"column:name"`
	Brand        string    `gorm:"column:brand"`
	Model        string    `gorm:"column:model"`
	year         string    `gorm:"column:year"`
	mileage      string    `gorm:"column:mileage"`
	transmission string    `gorm:"column:transmission"`
	fuel_type    string    `gorm:"column:fuel_type"`
	location     string    `gorm:"column:location"`
	description  string    `gorm:"column:description"`
	images       []string  `gorm:"column:images"`
	price        float64   `gorm:"column:price"`
	currency     string    `gorm:"column:currency"`
	IsSold       bool      `gorm:"column:is_sold"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *Cars) TableName() string {
	return "cars"
}
