package postgrestore

import (
	"e-wallet/domain/user"
	"time"
)

const (
	UsersTableName       = "users"
	UserProfilesTableName = "user_profiles"
)

type User struct {
	ID              string    `gorm:"type:uuid"`
	Username        string    `gorm:"uniqueIndex;not null"`
	Email           string    `gorm:"uniqueIndex;not null"`
	PasswordHash    string    `gorm:"not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	IsEmailVerified bool      `gorm:"default:false"`
}

type UserProfile struct {
	ID        int       `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	Name      string
	Email     string
	Avatar    string
	Phone     string
	IDNumber  string
	BirthYear int
	Gender    string
	Team      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain converts schema to domain model
func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:              u.ID,
		Username:        u.Username,
		Email:           u.Email,
		PasswordHash:    u.PasswordHash,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		IsEmailVerified: u.IsEmailVerified,
	}
}

// ToDomain converts schema to domain model
func (up *UserProfile) ToDomain() *user.UserProfile {
	return &user.UserProfile{
		ID:        up.ID,
		UserID:    up.UserID,
		Name:      up.Name,
		Email:     up.Email,
		Avatar:    up.Avatar,
		Phone:     up.Phone,
		IDNumber:  up.IDNumber,
		BirthYear: up.BirthYear,
		Gender:    up.Gender,
		Team:      up.Team,
	}
}