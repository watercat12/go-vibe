package user

import (
	"time"

	"e-wallet/pkg"
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

func NewProfile(userID, displayName, avatarURL, phoneNumber, nationalID, gender, team string, birthYear int) *Profile {
	return &Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      userID,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
		PhoneNumber: phoneNumber,
		NationalID:  nationalID,
		BirthYear:   birthYear,
		Gender:      gender,
		Team:        team,
	}
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