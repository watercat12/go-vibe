package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               string    `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	Username         string    `json:"username" gorm:"uniqueIndex;not null"`
	Email            string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string    `json:"-" gorm:"not null"`
	IsEmailVerified  bool      `json:"is_email_verified" gorm:"default:false"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	LoginUser(ctx context.Context, req *LoginUserRequest) (*User, error)
	UpdateProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*Profile, error)
	GetProfile(ctx context.Context, userID string) (*Profile, error)
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}