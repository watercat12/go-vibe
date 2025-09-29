package user

import (
	"context"
	"e-wallet/internal/domain/user"
	"e-wallet/internal/ports"
	"e-wallet/pkg"
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) user.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	hashedPassword, err := user.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	return s.repo.Create(ctx, u)
}
