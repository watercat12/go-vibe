package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashedpassword"

	user := NewUser(username, email, passwordHash)

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, passwordHash, user.PasswordHash)
	assert.False(t, user.IsEmailVerified)
}

func TestHashPassword(t *testing.T) {
	password := "password123"

	hashed, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)

	// Verify the hash works
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	assert.NoError(t, err)
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectError    bool
	}{
		{
			name:           "success - correct password",
			hashedPassword: string(hashed),
			password:       password,
			expectError:    false,
		},
		{
			name:           "error - incorrect password",
			hashedPassword: string(hashed),
			password:       "wrongpassword",
			expectError:    true,
		},
		{
			name:           "error - empty password",
			hashedPassword: string(hashed),
			password:       "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassword(tt.hashedPassword, tt.password)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}