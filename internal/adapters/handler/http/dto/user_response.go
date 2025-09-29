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

type UpdateProfileResponse struct {
	Profile *ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	PhoneNumber string `json:"phone_number"`
	NationalID  string `json:"national_id"`
	BirthYear   int    `json:"birth_year"`
	Gender      string `json:"gender"`
	Team        string `json:"team"`
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

func NewCreateUserResponse(token string) *CreateUserResponse {
	return &CreateUserResponse{
		Token: token,
	}
}

func NewLoginUserResponse(user *user.User, token string) *LoginUserResponse {
	return &LoginUserResponse{
		User:  NewUserResponse(user),
		Token: token,
	}
}

func NewUpdateProfileResponse(profile *user.Profile) *UpdateProfileResponse {
	return &UpdateProfileResponse{
		Profile: NewProfileResponse(profile),
	}
}

func NewProfileResponse(profile *user.Profile) *ProfileResponse {
	return &ProfileResponse{
		ID:          profile.ID,
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		PhoneNumber: profile.PhoneNumber,
		NationalID:  profile.NationalID,
		BirthYear:   profile.BirthYear,
		Gender:      profile.Gender,
		Team:        profile.Team,
	}
}
