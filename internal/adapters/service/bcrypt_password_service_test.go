package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordService_HashPassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "success - hash valid password",
			password:      "password123",
			expectedError: nil,
		},
		{
			name:          "success - hash empty password",
			password:      "",
			expectedError: nil,
		},
		{
			name:          "success - hash long password",
			password:      "thisisaverylongpasswordthatshouldstillworkwithbcrypt",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewBcryptPasswordService()

			hashed, err := service.HashPassword(tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Empty(t, hashed)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashed)
				assert.NotEqual(t, tt.password, hashed)

				// Verify the hash is valid bcrypt
				err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(tt.password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestBcryptPasswordService_CheckPassword(t *testing.T) {
	service := NewBcryptPasswordService()
	password := "password123"
	hashed, _ := service.HashPassword(password)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectedError  error
	}{
		{
			name:           "success - correct password",
			hashedPassword: hashed,
			password:       "password123",
			expectedError:  nil,
		},
		{
			name:           "error - incorrect password",
			hashedPassword: hashed,
			password:       "wrongpassword",
			expectedError:  bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:           "error - empty password",
			hashedPassword: hashed,
			password:       "",
			expectedError:  bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:           "error - invalid hash",
			hashedPassword: "invalidhash",
			password:       "password123",
			expectedError:  bcrypt.ErrHashTooShort, // or other bcrypt error
		},
		{
			name:           "error - empty hash",
			hashedPassword: "",
			password:       "password123",
			expectedError:  bcrypt.ErrHashTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CheckPassword(tt.hashedPassword, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}