package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username         string    `json:"username" gorm:"uniqueIndex;not null"`
	Email            string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash     string    `json:"-" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsEmailVerified  bool      `json:"is_email_verified" gorm:"default:false"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserProfile struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone"`
	IDNumber  string `json:"id_number"`
	BirthYear int    `json:"birth_year"`
	Gender    string `json:"gender"`
	Team      string `json:"team"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	UpdateProfile(ctx context.Context, profile *UserProfile) error
}

type UserService interface {
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