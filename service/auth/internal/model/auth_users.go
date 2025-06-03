package model

import "time"

type AuthUsers struct {
	ID             string    `gorm:"primaryKey;column:id"`
	Email          string    `gorm:"column:email"`
	HashedPassword string    `gorm:"column:hashed_password"`
	Provider       string    `gorm:"column:provider"`
	IsVerified     bool      `gorm:"column:is_verified"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *AuthUsers) TableName() string {
	return "auth_users"
}
