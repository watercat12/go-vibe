package postgres

import (
	"e-wallet/internal/domain/profile"
	"e-wallet/internal/domain/user"
	"time"
)

const (
	UsersTableName      = "users"
	UserProfilesTableName = "user_profiles"
)

type User struct {
	ID                  string
	Username            string
	Email               string
	PasswordHash        string
	IsEmailVerified     bool
	IsProfileCompleted  bool
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
}

type UserProfile struct {
	UserID      string
	DisplayName string
	AvatarURL   *string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:                 u.ID,
		Username:           u.Username,
		Email:              u.Email,
		PasswordHash:       u.PasswordHash,
		IsEmailVerified:    u.IsEmailVerified,
		IsProfileCompleted: u.IsProfileCompleted,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}

func (up *UserProfile) ToDomain() *profile.Profile {
	return &profile.Profile{
		UserID:      up.UserID,
		DisplayName: up.DisplayName,
		AvatarURL:   up.AvatarURL,
		PhoneNumber: up.PhoneNumber,
		NationalID:  up.NationalID,
		BirthYear:   up.BirthYear,
		Gender:      up.Gender,
		Team:        up.Team,
		CreatedAt:   up.CreatedAt,
		UpdatedAt:   up.UpdatedAt,
	}
}
