package postgres

import (
	"e-wallet/internal/domain/user"
	"time"
)

const (
	UsersTableName     = "users"
)

type User struct {
	ID              string    `gorm:"type:uuid"`
	Username        string    `gorm:"uniqueIndex;not null"`
	Email           string    `gorm:"uniqueIndex;not null"`
	PasswordHash    string    `gorm:"not null"`
	IsEmailVerified bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

// ToDomain converts schema to domain model
func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}