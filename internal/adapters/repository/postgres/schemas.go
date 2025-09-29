package postgres

import (
	"e-wallet/internal/domain/user"
	"time"
)

const (
	UsersTableName     = "users"
	ProfilesTableName  = "profiles"
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

type Profile struct {
	ID          string    `gorm:"type:uuid"`
	UserID      string    `gorm:"not null"`
	DisplayName string
	AvatarURL   string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Profile) TableName() string {
	return "profiles"
}

// ToDomain converts schema to domain model
func (p *Profile) ToDomain() *user.Profile {
	return &user.Profile{
		ID:          p.ID,
		UserID:      p.UserID,
		DisplayName: p.DisplayName,
		AvatarURL:   p.AvatarURL,
		PhoneNumber: p.PhoneNumber,
		NationalID:  p.NationalID,
		BirthYear:   p.BirthYear,
		Gender:      p.Gender,
		Team:        p.Team,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}