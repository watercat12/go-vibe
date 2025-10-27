package profile

import (
	"time"
)

type Profile struct {
	UserID       string
	DisplayName  string
	AvatarURL    *string
	PhoneNumber  string
	NationalID   string
	BirthYear    int
	Gender       string
	Team         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UpdateProfileRequest struct {
	DisplayName string
	AvatarURL   *string
	PhoneNumber string
	NationalID  string
	BirthYear   int
	Gender      string
	Team        string
}