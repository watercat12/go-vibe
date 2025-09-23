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
	ID        int       `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

type UserProfile struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"not null"`
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

func (UserProfile) TableName() string {
	return "user_profiles"
}

// ToDomain converts schema to domain model
func (u *User) ToDomain() *user.User {
	return &user.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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