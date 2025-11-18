package ports

import (
	"context"
	"e-wallet/internal/domain/profile"
)

type ProfileService interface {
	UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfileRequest) (*profile.Profile, error)
	GetProfile(ctx context.Context, userID string) (*profile.Profile, error)
}