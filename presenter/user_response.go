package presenter

import (
	"e-wallet/domain/user"
	"time"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewUserResponse(user *user.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewCreateUserResponse(token string) *CreateUserResponse {
	return &CreateUserResponse{
		Token: token,
	}
}

func NewLoginResponse(token string) *LoginResponse {
	return &LoginResponse{
		Token: token,
	}
}
