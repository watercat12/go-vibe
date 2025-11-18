package dto

import (
	"e-wallet/internal/domain/profile"
	"time"
)

type ProfileResponse struct {
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	AvatarURL   *string   `json:"avatar_url"`
	PhoneNumber string    `json:"phone_number"`
	NationalID  string    `json:"national_id"`
	BirthYear   int       `json:"birth_year"`
	Gender      string    `json:"gender"`
	Team        string    `json:"team"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProfileResponse(p *profile.Profile) *ProfileResponse {
	return &ProfileResponse{
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