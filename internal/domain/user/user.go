package user

import (
	"time"

	"e-wallet/pkg"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                  string
	Username            string
	Email               string
	PasswordHash        string
	IsEmailVerified     bool
	IsProfileCompleted  bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

type LoginUserRequest struct {
	Email    string
	Password string
}

func NewUser(username, email, passwordHash string) *User {
	return &User{
		ID:               pkg.NewUUIDV7(),
		Username:         username,
		Email:            email,
		PasswordHash:     passwordHash,
		IsProfileCompleted: false,
	}
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