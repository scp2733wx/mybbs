package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex:uk_users_username;size:32;not null" json:"username"`
	Name         string         `gorm:"size:32;not null" json:"name"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Role         string         `gorm:"type:enum('student','admin');size:32;default:student" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_users_deleted_at" json:"-"`
}

type JWT_user struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	User        *User  `json:"user"`
}
