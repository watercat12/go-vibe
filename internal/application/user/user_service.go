package user

import (
	"context"
	"e-wallet/internal/domain/user"
	"e-wallet/internal/ports"
	"e-wallet/pkg"
)

type userService struct {
	repo            ports.UserRepository
	profileRepo     ports.ProfileRepository
}

func NewUserService(repo ports.UserRepository, profileRepo ports.ProfileRepository) user.UserService {
	return &userService{repo: repo, profileRepo: profileRepo}
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

func (s *userService) LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.User, error) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := user.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req *user.UpdateProfileRequest) (*user.Profile, error) {
	// Upsert profile
	profile := &user.Profile{
		ID:          pkg.NewUUIDV7(),
		UserID:      userID,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		PhoneNumber: req.PhoneNumber,
		NationalID:  req.NationalID,
		BirthYear:   req.BirthYear,
		Gender:      req.Gender,
		Team:        req.Team,
	}

	p, err := s.profileRepo.Upsert(ctx, profile)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*user.Profile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}
