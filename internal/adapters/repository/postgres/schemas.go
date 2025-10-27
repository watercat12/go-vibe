package postgres

import (
	"e-wallet/internal/domain/user"
	"time"
)

const (
	UsersTableName    = "users"
)

type User struct {
	ID              string
	Username        string
	Email           string
	PasswordHash    string
	IsEmailVerified bool
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

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
