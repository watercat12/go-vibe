package ports

import (
	"context"
	"e-wallet/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
}