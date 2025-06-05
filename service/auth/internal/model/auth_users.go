package model

import "time"

type AuthUsers struct {
	Username       string    `gorm:"primaryKey;column:username"`
	FullName       string    `gorm:"column:full_name"`
	Email          string    `gorm:"column:email"`
	HashedPassword string    `gorm:"column:hashed_password"`
	Role           string    `gorm:"column:role"`
	IsVerified     bool      `gorm:"column:is_verified"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *AuthUsers) TableName() string {
	return "auth_users"
}
