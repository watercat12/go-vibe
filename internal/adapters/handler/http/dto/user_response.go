package dto

import (
	"e-wallet/internal/domain/user"
	"time"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

type LoginUserResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}


func NewUserResponse(user *user.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NewCreateUserResponse() *CreateUserResponse {
	return &CreateUserResponse{}
}

func NewLoginUserResponse(user *user.User, token string) *LoginUserResponse {
	return &LoginUserResponse{
		User:  NewUserResponse(user),
		Token: token,
	}
}

