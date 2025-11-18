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
		mockSetup     func(*mocks.MockUserRepository, *mocks.MockPasswordService)
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
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
				passwordService.EXPECT().HashPassword("password123").Return(hashPassword("password123"), nil).Once()
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
			name: "error - password hashing fails",
			request: &user.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
				passwordService.EXPECT().HashPassword("password123").Return("", errors.New("hash error")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("hash error"),
		},
		{
			name: "error - repository create fails",
			request: &user.CreateUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
				passwordService.EXPECT().HashPassword("password123").Return(hashPassword("password123"), nil).Once()
				userRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			passwordService := mocks.NewMockPasswordService(t)

			tt.mockSetup(userRepo, passwordService)

			service := NewUserService(userRepo, passwordService)
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
		mockSetup     func(*mocks.MockUserRepository, *mocks.MockPasswordService)
		expectedUser  *user.User
		expectedError error
	}{
		{
			name: "success - login user",
			request: &user.LoginUserRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
				userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(&user.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: hashedPwd,
				}, nil).Once()
				passwordService.EXPECT().CheckPassword(mock.Anything, "password123").Return(nil).Once()
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
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
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
			mockSetup: func(userRepo *mocks.MockUserRepository, passwordService *mocks.MockPasswordService) {
				userRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(&user.User{
					ID:           "user-123",
					Email:        "test@example.com",
					PasswordHash: hashedPwd,
				}, nil).Once()
				passwordService.EXPECT().CheckPassword(mock.Anything, "wrongpassword").Return(bcrypt.ErrMismatchedHashAndPassword).Once()
			},
			expectedUser:  nil,
			expectedError: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := mocks.NewMockUserRepository(t)
			passwordService := mocks.NewMockPasswordService(t)

			tt.mockSetup(userRepo, passwordService)

			service := NewUserService(userRepo, passwordService)
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