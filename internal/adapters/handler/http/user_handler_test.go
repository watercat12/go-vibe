package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/config"
	"e-wallet/internal/domain/user"
	"e-wallet/mocks"
	"e-wallet/pkg/logger"
)

func TestServer_CreateUser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      dto.CreateUserRequest
		mockSetup        func(*mocks.MockUserService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name: "success - create user",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					CreateUser(mock.Anything, mock.MatchedBy(func(req *user.CreateUserRequest) bool {
						return req.Username == "testuser"
					})).
					Return(&user.User{
						ID:           "user-123",
						Username:     "testuser",
						Email:        "test@example.com",
						PasswordHash: "hashed",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: &dto.CreateUserResponse{
					Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOnsidXNlcl9pZCI6InVzZXItMTIzIn0sImV4cCI6MTc1MzM5MzYwMH0.test", // approximate, will adjust
				},
			},
		},
		{
			name: "error - invalid request body",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as bind fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name: "error - validation fails",
			requestBody: dto.CreateUserRequest{
				Username: "",
				Email:    "invalid-email",
				Password: "password123",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as validation fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name: "error - service fails",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					CreateUser(mock.Anything, mock.MatchedBy(func(req *user.CreateUserRequest) bool {
						return req.Username == "testuser" && req.Email == "test@example.com" && req.Password == "Password123!"
					})).
					Return(nil, errors.New("service error")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: dto.InternalErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			userSvc := mocks.NewMockUserService(t)

			// Setup mocks
			tt.mockSetup(userSvc)

			// Create server
			e := echo.New()
			var reqBody bytes.Buffer
			if tt.name != "error - invalid request body" {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			} else {
				reqBody.WriteString("invalid json")
			}
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", &reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set validator for tests that reach validation
			if tt.name != "error - invalid request body" {
				v := validator.New()
				dto.RegisterCustomValidations(v)
				e.Validator = &CustomValidator{validator: v}
			}

			// Create server instance
			s := &Server{
				UserService: userSvc,
				Config: &config.Config{
					JWTSecret: "test-secret",
				},
				Logger: logger.NOOPLogger,
			}

			// Execute
			err := s.CreateUser(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.name == "success - create user" {
				// For success, check that data is present and token is string
				var resp dto.Response
				json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.Equal(t, http.StatusOK, resp.Status)
				assert.Equal(t, "OK", resp.Message)
				data, ok := resp.Data.(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, data["token"])
			} else {
				expectedJSON, _ := json.Marshal(tt.expectedResponse)
				assert.JSONEq(t, string(expectedJSON), rec.Body.String())
			}
		})
	}
}

func TestServer_LoginUser(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      dto.LoginUserRequest
		mockSetup        func(*mocks.MockUserService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name: "success - login user",
			requestBody: dto.LoginUserRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					LoginUser(mock.Anything, mock.MatchedBy(func(req *user.LoginUserRequest) bool {
						return req.Email == "test@example.com" && req.Password == "Password123!"
					})).
					Return(&user.User{
						ID:           "user-123",
						Username:     "testuser",
						Email:        "test@example.com",
						PasswordHash: "hashed",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: &dto.LoginUserResponse{
					User: &dto.UserResponse{
						ID:       "user-123",
						Username: "testuser",
						Email:    "test@example.com",
					},
					Token: "token",
				},
			},
		},
		{
			name: "error - invalid request body",
			requestBody: dto.LoginUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as bind fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name: "error - validation fails",
			requestBody: dto.LoginUserRequest{
				Email:    "invalid-email",
				Password: "",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as validation fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name: "error - service fails",
			requestBody: dto.LoginUserRequest{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					LoginUser(mock.Anything, mock.MatchedBy(func(req *user.LoginUserRequest) bool {
						return req.Email == "test@example.com" && req.Password == "Password123!"
					})).
					Return(nil, errors.New("service error")).
					Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			userSvc := mocks.NewMockUserService(t)

			// Setup mocks
			tt.mockSetup(userSvc)

			// Create server
			e := echo.New()
			var reqBody bytes.Buffer
			if tt.name != "error - invalid request body" {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			} else {
				reqBody.WriteString("invalid json")
			}
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", &reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set validator for tests that reach validation
			if tt.name != "error - invalid request body" {
				v := validator.New()
				dto.RegisterCustomValidations(v)
				e.Validator = &CustomValidator{validator: v}
			}

			// Create server instance
			s := &Server{
				UserService: userSvc,
				Config: &config.Config{
					JWTSecret: "test-secret",
				},
				Logger: logger.NOOPLogger,
			}

			// Execute
			err := s.LoginUser(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.name == "success - login user" {
				// For success, check that data is present
				var resp dto.Response
				json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.Equal(t, http.StatusOK, resp.Status)
				assert.Equal(t, "OK", resp.Message)
				data, ok := resp.Data.(map[string]any)
				assert.True(t, ok)
				assert.NotNil(t, data["user"])
				assert.NotEmpty(t, data["token"])
			} else {
				expectedJSON, _ := json.Marshal(tt.expectedResponse)
				assert.JSONEq(t, string(expectedJSON), rec.Body.String())
			}
		})
	}
}

func TestServer_UpdateProfile(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		requestBody      dto.UpdateProfileRequest
		mockSetup        func(*mocks.MockUserService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name:   "success - update profile",
			userID: "user-123",
			requestBody: dto.UpdateProfileRequest{
				Username:    "newusername",
				DisplayName: "New Display Name",
				AvatarURL:   "http://example.com/avatar.jpg",
				PhoneNumber: "1234567890",
				NationalID:  "123456789",
				BirthYear:   1990,
				Gender:      "male",
				Team:        "Back End",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					UpdateProfile(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(req *user.UpdateProfileRequest) bool {
						return req.Username == "newusername" && req.DisplayName == "New Display Name"
					})).
					Return(&user.Profile{
						ID:          "profile-123",
						UserID:      "user-123",
						DisplayName: "New Display Name",
						AvatarURL:   "http://example.com/avatar.jpg",
						PhoneNumber: "1234567890",
						NationalID:  "123456789",
						BirthYear:   1990,
						Gender:      "male",
						Team:        "Back End",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: &dto.UpdateProfileResponse{
					Profile: &dto.ProfileResponse{
						ID:          "profile-123",
						UserID:      "user-123",
						DisplayName: "New Display Name",
						AvatarURL:   "http://example.com/avatar.jpg",
						PhoneNumber: "1234567890",
						NationalID:  "123456789",
						BirthYear:   1990,
						Gender:      "male",
						Team:        "Back End",
					},
				},
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			requestBody: dto.UpdateProfileRequest{
				Username:    "newusername",
				DisplayName: "New Display Name",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
		{
			name:   "error - invalid request body",
			userID: "user-123",
			requestBody: dto.UpdateProfileRequest{
				Username:    "newusername",
				DisplayName: "New Display Name",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as bind fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - validation fails",
			userID: "user-123",
			requestBody: dto.UpdateProfileRequest{
				Username: "",
				Team:     "invalid-team",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				// No mock setup needed as validation fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - service fails",
			userID: "user-123",
			requestBody: dto.UpdateProfileRequest{
				Username:    "newusername",
				DisplayName: "New Display Name",
				Team:        "Back End",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					UpdateProfile(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(req *user.UpdateProfileRequest) bool {
						return req.Username == "newusername" && req.DisplayName == "New Display Name" && req.Team == "Back End"
					})).
					Return(nil, errors.New("service error")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: dto.InternalErrorResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			userSvc := mocks.NewMockUserService(t)

			// Setup mocks
			tt.mockSetup(userSvc)

			// Create server
			e := echo.New()
			var reqBody bytes.Buffer
			if tt.name != "error - invalid request body" {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			} else {
				reqBody.WriteString("invalid json")
			}
			req := httptest.NewRequest(http.MethodPut, "/api/user/profile", &reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set validator for tests that reach validation
			if tt.name != "error - unauthorized" && tt.name != "error - invalid request body" {
				v := validator.New()
				dto.RegisterCustomValidations(v)
				e.Validator = &CustomValidator{validator: v}
			}

			// Set user ID if provided
			if tt.userID != "" {
				c.Set(UserIDKey, tt.userID)
			}

			// Create server instance
			s := &Server{
				UserService: userSvc,
				Config: &config.Config{
					JWTSecret: "test-secret",
				},
				Logger: logger.NOOPLogger,
			}

			// Execute
			err := s.UpdateProfile(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}