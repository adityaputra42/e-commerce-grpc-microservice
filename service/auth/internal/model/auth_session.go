package model

import "time"

type AuthSessions struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Username     string    `gorm:"primaryKey;column:username"`
	RefreshToken string    `gorm:"column:refresh_token"`
	UserAgent    string    `gorm:"column:user_agent"`
	ClientIp     string    `gorm:"column:client_ip"`
	IsBlocked    bool      `gorm:"column:is_blocked"`
	ExpiredAt    time.Time `gorm:"column:expired_at"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (u *AuthSessions) TableName() string {
	return "auth_session"
}
