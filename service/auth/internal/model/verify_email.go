package model

import "time"

type VerifyEmail struct {
	ID         int64     `gorm:"primaryKey;column:id"`
	Username   string    `gorm:"primaryKey;column:username"`
	Email      string    `gorm:"column:email"`
	SecretCode string    `gorm:"column:secret_code"`
	IsUsed     bool      `gorm:"column:is_used"`
	ExpiredAt  time.Time `gorm:"column:expired_at"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *VerifyEmail) TableName() string {
	return "verify_email"
}
