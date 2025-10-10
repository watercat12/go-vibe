package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestContainFirst(t *testing.T) {
	tests := []struct {
		name     string
		elems    []string
		v        string
		expected bool
	}{
		{
			name:     "match first element",
			elems:    []string{"/api", "/auth"},
			v:        "/api/users",
			expected: true,
		},
		{
			name:     "match second element",
			elems:    []string{"/api", "/auth"},
			v:        "/auth/login",
			expected: true,
		},
		{
			name:     "no match",
			elems:    []string{"/api", "/auth"},
			v:        "/health",
			expected: false,
		},
		{
			name:     "empty slice",
			elems:    []string{},
			v:        "/api",
			expected: false,
		},
		{
			name:     "exact match",
			elems:    []string{"/api"},
			v:        "/api",
			expected: true,
		},
		{
			name:     "no prefix match",
			elems:    []string{"/api"},
			v:        "/testapi",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainFirst(tt.elems, tt.v)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewAuthentication(t *testing.T) {
	keyLookup := "header:Authorization"
	authScheme := "Bearer"
	skipperPath := []string{"/health", "/auth"}

	auth := NewAuthentication(keyLookup, authScheme, skipperPath)

	assert.NotNil(t, auth)
	assert.Equal(t, skipperPath, auth.SkipperPath)
	assert.Equal(t, keyLookup, auth.KeyLookup)
	assert.Equal(t, authScheme, auth.AuthScheme)
}

func TestAuthentication_ValidateAccessToken(t *testing.T) {
	// Set up environment
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	defer os.Setenv("JWT_SECRET_KEY", originalSecret)
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	auth := NewAuthentication("header:Authorization", "Bearer", []string{})

	tests := []struct {
		name        string
		token       string
		setupToken  func() string
		expected    bool
		expectError bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := CreateAccessToken(DefaultExpiredTime, TokenPayload{UserID: "user-123"}, "test-secret-key")
				return token
			},
			expected:    true,
			expectError: false,
		},
		{
			name:        "empty token",
			token:       "",
			expected:    false,
			expectError: true,
		},
		{
			name: "expired token",
			setupToken: func() string {
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["sub"] = TokenPayload{UserID: "user-123"}
				claims["exp"] = float64(0) // expired
				tokenString, _ := token.SignedString([]byte("test-secret-key"))
				return tokenString
			},
			expected:    false,
			expectError: true,
		},
		{
			name: "invalid token",
			token: "invalid.jwt.token",
			expected: false,
			expectError: true,
		},
		{
			name: "token with empty user_id",
			setupToken: func() string {
				token, _ := CreateAccessToken(DefaultExpiredTime, TokenPayload{UserID: ""}, "test-secret-key")
				return token
			},
			expected:    false,
			expectError: true,
		},
		{
			name: "token with invalid sub claim",
			setupToken: func() string {
				token := jwt.New(jwt.SigningMethodHS256)
				claims := token.Claims.(jwt.MapClaims)
				claims["exp"] = time.Now().UTC().Add(DefaultExpiredTime).Unix()
				claims["sub"] = "invalid json"
				tokenString, _ := token.SignedString([]byte("test-secret-key"))
				return tokenString
			},
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			if tt.setupToken != nil {
				token = tt.setupToken()
			} else {
				token = tt.token
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			result, err := auth.ValidateAccessToken(token, c)

			assert.Equal(t, tt.expected, result)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthentication_Middleware(t *testing.T) {
	// Set up environment
	originalSecret := os.Getenv("JWT_SECRET_KEY")
	defer os.Setenv("JWT_SECRET_KEY", originalSecret)
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	auth := NewAuthentication("header:Authorization", "Bearer", []string{"/health"})

	e := echo.New()
	e.Use(auth.Middleware())
	e.Use(CheckUserTypeMiddleware())

	// Handler that just returns OK
	handler := func(c echo.Context) error {
		userID, ok := c.Get(UserIDKey).(string)
		if !ok {
			return c.String(http.StatusOK, "no user")
		}
		return c.String(http.StatusOK, "user: "+userID)
	}

	e.GET("/protected", handler)
	e.GET("/health", handler)

	tests := []struct {
		name           string
		path           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "skip path - no auth required",
			path:           "/health",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedBody:   "no user",
		},
		{
			name: "protected path - valid token",
			path: "/protected",
			authHeader: func() string {
				token, _ := CreateAccessToken(DefaultExpiredTime, TokenPayload{UserID: "user-123"}, "test-secret-key")
				return "Bearer " + token
			}(),
			expectedStatus: http.StatusOK,
			expectedBody:   "user: user-123",
		},
		{
			name:           "protected path - no token",
			path:           "/protected",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
		},
		{
			name:           "protected path - invalid token",
			path:           "/protected",
			authHeader:     "Bearer invalid.token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req.Header.Set("Authorization", tt.authHeader)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rec.Body.String())
			}
		})
	}
}

func TestCheckUserTypeMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(CheckUserTypeMiddleware())

	handler := func(c echo.Context) error {
		userID, ok := c.Get(UserIDKey).(string)
		if !ok {
			return c.String(http.StatusOK, "no user")
		}
		return c.String(http.StatusOK, "user: "+userID)
	}

	e.GET("/test", handler)

	tests := []struct {
		name         string
		setUserClaim bool
		userID       string
		expectedBody string
	}{
		{
			name:         "user claim set",
			setUserClaim: true,
			userID:       "user-123",
			expectedBody: "user: user-123",
		},
		{
			name:         "no user claim",
			setUserClaim: false,
			expectedBody: "no user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.setUserClaim {
				c.Set(UserClaimKey, &TokenPayload{UserID: tt.userID})
			}

			err := CheckUserTypeMiddleware()(func(c echo.Context) error {
				return handler(c)
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}
}