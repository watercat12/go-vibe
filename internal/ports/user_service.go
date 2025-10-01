package ports

import (
	"context"

	"e-wallet/internal/domain/user"
)

type UserService interface {
	CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error)
	LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.User, error)
	UpdateProfile(ctx context.Context, userID string, req *user.UpdateProfileRequest) (*user.Profile, error)
	GetProfile(ctx context.Context, userID string) (*user.Profile, error)
}