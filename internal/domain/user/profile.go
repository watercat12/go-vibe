package user

import (
	"time"
)

type Profile struct {
	ID          string
	UserID      string
	DisplayName string
	AvatarURL   string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateProfileRequest struct {
	Username    string
	DisplayName string
	AvatarURL   string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
}