package http

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessToken(t *testing.T) {
	tests := []struct {
		name           string
		ttl            time.Duration
		payload        TokenPayload
		secretJWTKey   string
		expectedError  bool
	}{
		{
			name: "success - create access token",
			ttl:  DefaultExpiredTime,
			payload: TokenPayload{
				UserID: "user-123",
			},
			secretJWTKey:  "secret",
			expectedError: false,
		},
		{
			name: "success - create access token with custom ttl",
			ttl:  time.Hour,
			payload: TokenPayload{
				UserID: "user-456",
			},
			secretJWTKey:  "another-secret",
			expectedError: false,
		},
		{
			name: "success - empty secret key",
			ttl:  DefaultExpiredTime,
			payload: TokenPayload{
				UserID: "user-123",
			},
			secretJWTKey:  "",
			expectedError: false, // JWT allows empty secret
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateAccessToken(tt.ttl, tt.payload, tt.secretJWTKey)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Validate the token
				claims, err := ValidateToken(token, tt.secretJWTKey)
				assert.NoError(t, err)
				assert.NotNil(t, claims)

				// Decode the payload
				payload, err := DecodeToken(claims)
				assert.NoError(t, err)
				assert.Equal(t, tt.payload.UserID, payload.UserID)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret"
	validPayload := TokenPayload{UserID: "user-123"}
	validToken, _ := CreateAccessToken(DefaultExpiredTime, validPayload, secret)

	expiredToken, _ := CreateAccessToken(-time.Hour, validPayload, secret) // Expired

	tests := []struct {
		name           string
		token          string
		secretJWTKey   string
		expectedError  bool
	}{
		{
			name:          "success - validate valid token",
			token:         validToken,
			secretJWTKey:  secret,
			expectedError: false,
		},
		{
			name:          "error - invalid token",
			token:         "invalid.token.here",
			secretJWTKey:  secret,
			expectedError: true,
		},
		{
			name:          "error - wrong secret",
			token:         validToken,
			secretJWTKey:  "wrong-secret",
			expectedError: true,
		},
		{
			name:          "error - expired token",
			token:         expiredToken,
			secretJWTKey:  secret,
			expectedError: true,
		},
		{
			name:          "error - invalid claims type",
			token:         func() string {
				token := jwt.New(jwt.SigningMethodHS256)
				token.Claims = jwt.MapClaims{}
				tokenString, _ := token.SignedString([]byte(secret))
				// Manually corrupt the token to have invalid claims
				return tokenString[:len(tokenString)-5] + "xxxxx"
			}(),
			secretJWTKey:  secret,
			expectedError: true,
		},
		{
			name:          "error - token not valid",
			token:         func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": map[string]interface{}{"user_id": "user-123"},
					"exp": time.Now().Add(-time.Hour).Unix(), // Expired
				})
				tokenString, _ := token.SignedString([]byte(secret))
				return tokenString
			}(),
			secretJWTKey:  secret,
			expectedError: true,
		},
		{
			name:          "error - empty token",
			token:         "",
			secretJWTKey:  secret,
			expectedError: true,
		},
		{
			name: "error - wrong signing method",
			token: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": map[string]interface{}{"user_id": "user-123"},
					"exp": time.Now().Add(time.Hour).Unix(),
				})
				token.Header["alg"] = "RS256"
				tokenString, _ := token.SignedString([]byte(secret))
				return tokenString
			}(),
			secretJWTKey:  secret,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token, tt.secretJWTKey)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestDecodeToken(t *testing.T) {
	validClaims := jwt.MapClaims{
		"sub": map[string]interface{}{
			"user_id": "user-123",
		},
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	invalidClaimsMissingSub := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	invalidClaimsInvalidSub := jwt.MapClaims{
		"sub": "invalid-json",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	invalidClaimsMarshalError := jwt.MapClaims{
		"sub": make(chan int), // Cannot marshal channel
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	tests := []struct {
		name           string
		claims         jwt.MapClaims
		expectedError  bool
		expectedUserID string
	}{
		{
			name:           "success - decode valid claims",
			claims:         validClaims,
			expectedError:  false,
			expectedUserID: "user-123",
		},
		{
			name:          "error - missing sub",
			claims:        invalidClaimsMissingSub,
			expectedError: true,
		},
		{
			name:          "error - invalid sub json",
			claims:        invalidClaimsInvalidSub,
			expectedError: true,
		},
		{
			name:          "error - marshal error",
			claims:        invalidClaimsMarshalError,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := DecodeToken(tt.claims)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, payload)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payload)
				assert.Equal(t, tt.expectedUserID, payload.UserID)
			}
		})
	}
}