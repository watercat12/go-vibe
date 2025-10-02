package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"e-wallet/internal/domain/user"
	"e-wallet/mocks"
)

// Helper function to hash password for tests
func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *user.CreateUserRequest
		mockSetup     func(*mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedUser  *user.User
		expectedError error
	}{
		{
			name: "success - create user",
			request: &user.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *user.User) bool {
					return u.Username == "testuser" && u.Email == "test@example.com" && len(u.PasswordHash) > 0
				})).Return(&user.User{
					ID:           "user-123",
					Username:     "testuser",
					Email:        "test@example.com",
					PasswordHash: hashPassword("password123"),
				}, nil).Once()
			},
			expectedUser: &user.User{
				ID:       "user-123",
				Username: "testuser",
				Email:    "test@example.com",
			},
			expectedError: nil,
		},
		{
			name: "error - repository create fails",
			request: &user.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			tt.mockSetup(userRepo, profileRepo)

			service := NewUserService(userRepo, profileRepo)
			result, err := service.CreateUser(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Username, result.Username)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	hashedPwd := hashPassword("password123")

	tests := []struct {
		name          string
		request       *user.LoginUserRequest
		mockSetup     func(*mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedUser  *user.User
		expectedError error
	}{
		{
			name: "success - login user",
			request: &user.LoginUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(&user.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: hashedPwd,
				}, nil).Once()
			},
			expectedUser: &user.User{
				ID:    "user-123",
				Email: "test@example.com",
			},
			expectedError: nil,
		},
		{
			name: "error - user not found",
			request: &user.LoginUserRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.EXPECT().GetByEmail(mock.Anything, "nonexistent@example.com").Return(nil, errors.New("user not found")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name: "error - incorrect password",
			request: &user.LoginUserRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(&user.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: hashedPwd,
				}, nil).Once()
			},
			expectedUser:  nil,
			expectedError: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			tt.mockSetup(userRepo, profileRepo)

			service := NewUserService(userRepo, profileRepo)
			result, err := service.LoginUser(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
			}
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		request         *user.UpdateProfileRequest
		mockSetup       func(*mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedProfile *user.Profile
		expectedError   error
	}{
		{
			name:   "success - update profile",
			userID: "user-123",
			request: &user.UpdateProfileRequest{
				DisplayName: "Test User",
				AvatarURL:   "http://example.com/avatar.jpg",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				profileRepo.EXPECT().Upsert(mock.Anything, mock.MatchedBy(func(p *user.Profile) bool {
					return p.UserID == "user-123" && p.DisplayName == "Test User"
				})).Return(&user.Profile{
					ID:          "profile-456",
					UserID:      "user-123",
					DisplayName: "Test User",
					AvatarURL:   "http://example.com/avatar.jpg",
				}, nil).Once()
			},
			expectedProfile: &user.Profile{
				ID:          "profile-456",
				UserID:      "user-123",
				DisplayName: "Test User",
				AvatarURL:   "http://example.com/avatar.jpg",
			},
			expectedError: nil,
		},
		{
			name:   "error - upsert fails",
			userID: "user-123",
			request: &user.UpdateProfileRequest{
				DisplayName: "Test User",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				profileRepo.EXPECT().Upsert(mock.Anything, mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expectedProfile: nil,
			expectedError:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			tt.mockSetup(userRepo, profileRepo)

			service := NewUserService(userRepo, profileRepo)
			result, err := service.UpdateProfile(context.Background(), tt.userID, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedProfile.UserID, result.UserID)
				assert.Equal(t, tt.expectedProfile.DisplayName, result.DisplayName)
			}
		})
	}
}

func TestUserService_GetProfile(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		mockSetup       func(*mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedProfile *user.Profile
		expectedError   error
	}{
		{
			name:   "success - get profile",
			userID: "user-123",
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				profileRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(&user.Profile{
					ID:     "profile-456",
					UserID: "user-123",
				}, nil).Once()
			},
			expectedProfile: &user.Profile{
				ID:     "profile-456",
				UserID: "user-123",
			},
			expectedError: nil,
		},
		{
			name:   "error - profile not found",
			userID: "user-123",
			mockSetup: func(userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				profileRepo.EXPECT().GetByUserID(mock.Anything, "user-123").Return(nil, errors.New("profile not found")).Once()
			},
			expectedProfile: nil,
			expectedError:   errors.New("profile not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			tt.mockSetup(userRepo, profileRepo)

			service := NewUserService(userRepo, profileRepo)
			result, err := service.GetProfile(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedProfile.UserID, result.UserID)
			}
		})
	}
}