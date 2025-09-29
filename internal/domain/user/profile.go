package user

import (
	"time"
)

type Profile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	PhoneNumber string    `json:"phone_number"`
	NationalID  string    `json:"national_id"`
	BirthYear   int       `json:"birth_year"`
	Gender      string    `json:"gender"`
	Team        string    `json:"team"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	NationalID  string `json:"national_id"`
	BirthYear   int    `json:"birth_year"`
	Gender      string `json:"gender"`
	Team        string `json:"team"`
}