package ports

import (
	"context"
	"e-wallet/internal/domain/profile"
)

type ProfileRepository interface {
	GetByUserID(ctx context.Context, userID string) (*profile.Profile, error)
	Upsert(ctx context.Context, profile *profile.Profile) (*profile.Profile, error)
	CheckNationalIDExists(ctx context.Context, nationalID string, excludeUserID string) (bool, error)
}