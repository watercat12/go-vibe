package ports

import (
	"context"
	"e-wallet/internal/domain/user"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *user.Profile) (*user.Profile, error)
	GetByUserID(ctx context.Context, userID string) (*user.Profile, error)
	Update(ctx context.Context, profile *user.Profile) (*user.Profile, error)
	Upsert(ctx context.Context, profile *user.Profile) (*user.Profile, error)
}