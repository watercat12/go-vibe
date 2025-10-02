package user

import (
	"context"
	"e-wallet/internal/domain/user"
	"e-wallet/internal/ports"
)

type userService struct {
	repo            ports.UserRepository
	profileRepo     ports.ProfileRepository
	passwordService ports.PasswordService
}

func NewUserService(repo ports.UserRepository, profileRepo ports.ProfileRepository, passwordService ports.PasswordService) ports.UserService {
	return &userService{repo: repo, profileRepo: profileRepo, passwordService: passwordService}
}

func (s *userService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := user.NewUser(req.Username, req.Email, hashedPassword)

	return s.repo.Create(ctx, u)
}

func (s *userService) LoginUser(ctx context.Context, req *user.LoginUserRequest) (*user.User, error) {
	u, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := s.passwordService.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID string, req *user.UpdateProfileRequest) (*user.Profile, error) {
	profile := user.NewProfile(userID, req.DisplayName, req.AvatarURL, req.PhoneNumber, req.NationalID, req.Gender, req.Team, req.BirthYear)

	p, err := s.profileRepo.Upsert(ctx, profile)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*user.Profile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}
