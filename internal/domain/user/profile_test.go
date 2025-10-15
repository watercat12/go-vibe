package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProfile(t *testing.T) {
	userID := "user-123"
	displayName := "Test User"
	avatarURL := "http://example.com/avatar.jpg"
	phoneNumber := "123456789"
	nationalID := "123456789012"
	gender := "male"
	team := "team1"
	birthYear := 1990

	profile := NewProfile(userID, displayName, avatarURL, phoneNumber, nationalID, gender, team, birthYear)

	assert.NotEmpty(t, profile.ID)
	assert.Equal(t, userID, profile.UserID)
	assert.Equal(t, displayName, profile.DisplayName)
	assert.Equal(t, avatarURL, profile.AvatarURL)
	assert.Equal(t, phoneNumber, profile.PhoneNumber)
	assert.Equal(t, nationalID, profile.NationalID)
	assert.Equal(t, birthYear, profile.BirthYear)
	assert.Equal(t, gender, profile.Gender)
	assert.Equal(t, team, profile.Team)
}