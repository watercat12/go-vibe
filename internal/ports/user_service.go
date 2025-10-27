package ports

import (
	"context"

	"e-wallet/internal/domain/user"
)

type UserService interface {
	CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error)
	LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.User, error)
}