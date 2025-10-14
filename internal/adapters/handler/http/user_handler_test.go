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
	"strings"
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
				Password: "TestPass123@!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					CreateUser(mock.Anything, mock.MatchedBy(func(req *user.CreateUserRequest) bool {
						return req.Username == "testuser" && req.Email == "test@example.com" && req.Password == "TestPass123@!"
					})).
					Return(&user.User{
						ID:              "user-123",
						Username:        "testuser",
						Email:           "test@example.com",
						PasswordHash:    "hashedpassword",
						IsEmailVerified: false,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    dto.NewCreateUserResponse(),
			},
		},
		{
			name: "error - invalid request body",
			requestBody: dto.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "TestPass123@!",
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
				Password: "short",
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
				Password: "TestPass123@!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					CreateUser(mock.Anything, mock.MatchedBy(func(req *user.CreateUserRequest) bool {
						return req.Username == "testuser" && req.Email == "test@example.com" && req.Password == "TestPass123@!"
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
			req := httptest.NewRequest(http.MethodPost, "/api/users", &reqBody)
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
				Logger:      logger.NOOPLogger,
			}

			// Execute
			err := s.CreateUser(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse dto.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			dataActualJson, _ := json.Marshal(actualResponse.Data)
			var dataActual dto.CreateUserResponse
			_ = json.Unmarshal(dataActualJson, &dataActual)

			assert.Equal(t, dto.Response{
				Status:  tt.expectedResponse.Status,
				Message: tt.expectedResponse.Message,
			}, dto.Response{
				Status:  actualResponse.Status,
				Message: actualResponse.Message,
			})
			if strings.Contains(tt.name, "success") {
				assert.Equal(t,dataActual,dto.CreateUserResponse{})
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
				Password: "TestPass123@!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					LoginUser(mock.Anything, mock.MatchedBy(func(req *user.LoginUserRequest) bool {
						return req.Email == "test@example.com" && req.Password == "TestPass123@!"
					})).
					Return(&user.User{
						ID:              "user-123",
						Username:        "testuser",
						Email:           "test@example.com",
						PasswordHash:    "hashedpassword",
						IsEmailVerified: true,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: dto.NewLoginUserResponse(&user.User{
					ID:              "user-123",
					Username:        "testuser",
					Email:           "test@example.com",
					PasswordHash:    "hashedpassword",
					IsEmailVerified: true,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}, "jwt-token"),
			},
		},
		{
			name: "error - invalid request body",
			requestBody: dto.LoginUserRequest{
				Email:    "test@example.com",
				Password: "TestPass123@!",
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
				Password: "TestPass123@!",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					LoginUser(mock.Anything, mock.MatchedBy(func(req *user.LoginUserRequest) bool {
						return req.Email == "test@example.com" && req.Password == "TestPass123@!"
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
			req := httptest.NewRequest(http.MethodPost, "/api/users/login", &reqBody)
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
				Logger:      logger.NOOPLogger,
				Config: &config.Config{
					JWTSecret: "test-secret",
				},
			}

			// Execute
			err := s.LoginUser(c)
			assert.NoError(t, err)

			// Assert
			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse dto.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			dataActualJson, _ := json.Marshal(actualResponse.Data)
			var dataActual dto.LoginUserResponse
			_ = json.Unmarshal(dataActualJson, &dataActual)

			assert.Equal(t, dto.Response{
				Status:  tt.expectedResponse.Status,
				Message: tt.expectedResponse.Message,
			}, dto.Response{
				Status:  actualResponse.Status,
				Message: actualResponse.Message,
			})
			if strings.Contains(tt.name, "success") {
				assert.NotEmpty(t,dataActual.Token)
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
				Username:    "updateduser",
				DisplayName: "Updated User",
				AvatarURL:   "https://example.com/avatar.jpg",
				PhoneNumber: "1234567890",
				NationalID:  "123456789",
				BirthYear:   1990,
				Gender:      "Male",
				Team:        "Back End",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					UpdateProfile(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(req *user.UpdateProfileRequest) bool {
						return req.Username == "updateduser" && req.DisplayName == "Updated User" && req.Team == "Back End"
					})).
					Return(&user.Profile{
						ID:          "profile-123",
						UserID:      "user-123",
						DisplayName: "Updated User",
						AvatarURL:   "https://example.com/avatar.jpg",
						PhoneNumber: "1234567890",
						NationalID:  "123456789",
						BirthYear:   1990,
						Gender:      "Male",
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
				Data: dto.NewUpdateProfileResponse(&user.Profile{
					ID:          "profile-123",
					UserID:      "user-123",
					DisplayName: "Updated User",
					AvatarURL:   "https://example.com/avatar.jpg",
					PhoneNumber: "1234567890",
					NationalID:  "123456789",
					BirthYear:   1990,
					Gender:      "Male",
					Team:        "Back End",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}),
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			requestBody: dto.UpdateProfileRequest{
				Username:    "updateduser",
				DisplayName: "Updated User",
				Team:        "Back End",
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
				Username:    "updateduser",
				DisplayName: "Updated User",
				Team:        "Back End",
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
				Team:     "Invalid Team",
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
				Username:    "updateduser",
				DisplayName: "Updated User",
				Team:        "Back End",
			},
			mockSetup: func(userSvc *mocks.MockUserService) {
				userSvc.EXPECT().
					UpdateProfile(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(req *user.UpdateProfileRequest) bool {
						return req.Username == "updateduser" && req.DisplayName == "Updated User" && req.Team == "Back End"
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
			req := httptest.NewRequest(http.MethodPut, "/api/users/profile", &reqBody)
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
				Logger:      logger.NOOPLogger,
			}

			// Execute
			err := s.UpdateProfile(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse dto.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			dataActualJson, _ := json.Marshal(actualResponse.Data)
			var dataActual dto.UpdateProfileResponse
			_ = json.Unmarshal(dataActualJson, &dataActual)

			assert.Equal(t, dto.Response{
				Status:  tt.expectedResponse.Status,
				Message: tt.expectedResponse.Message,
			}, dto.Response{
				Status:  actualResponse.Status,
				Message: actualResponse.Message,
			})
			if strings.Contains(tt.name, "success") {
				assert.NotNil(t, dataActual.Profile)
				assert.Equal(t, "profile-123", dataActual.Profile.ID)
				assert.Equal(t, "user-123", dataActual.Profile.UserID)
				assert.Equal(t, "Updated User", dataActual.Profile.DisplayName)
				assert.Equal(t, "Back End", dataActual.Profile.Team)
			}
		})
	}
}