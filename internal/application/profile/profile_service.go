package profile

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"e-wallet/internal/domain/profile"
	"e-wallet/internal/ports"
)

type profileService struct {
	userRepo    ports.UserRepository
	profileRepo ports.ProfileRepository
}

func NewProfileService(userRepo ports.UserRepository, profileRepo ports.ProfileRepository) ports.ProfileService {
	return &profileService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

func (s *profileService) UpdateProfile(ctx context.Context, userID string, req *profile.UpdateProfileRequest) (*profile.Profile, error) {
	// Validate input
	if err := s.validateUpdateProfileRequest(req); err != nil {
		return nil, err
	}

	// Check if national_id already exists for another user
	exists, err := s.profileRepo.CheckNationalIDExists(ctx, req.NationalID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("national ID already registered for another user")
	}

	// Upsert profile
	newProfile := &profile.Profile{
		UserID:      userID,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		PhoneNumber: req.PhoneNumber,
		NationalID:  req.NationalID,
		BirthYear:   req.BirthYear,
		Gender:      req.Gender,
		Team:        req.Team,
	}
	updatedProfile, err := s.profileRepo.Upsert(ctx, newProfile)
	if err != nil {
		return nil, err
	}

	// Mark profile as completed
	if err := s.userRepo.UpdateProfileCompleted(ctx, userID, true); err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

func (s *profileService) GetProfile(ctx context.Context, userID string) (*profile.Profile, error) {
	return s.profileRepo.GetByUserID(ctx, userID)
}

func (s *profileService) validateUpdateProfileRequest(req *profile.UpdateProfileRequest) error {
	if req.DisplayName == "" {
		return errors.New("display name is required")
	}
	if req.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if req.NationalID == "" {
		return errors.New("national ID is required")
	}
	if req.BirthYear <= 0 {
		return errors.New("birth year is required")
	}
	if req.Gender == "" {
		return errors.New("gender is required")
	}
	if req.Team == "" {
		return errors.New("team is required")
	}

	// Validate phone number format (basic validation)
	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	if !phoneRegex.MatchString(req.PhoneNumber) {
		return errors.New("invalid phone number format")
	}

	// Validate birth year
	currentYear := time.Now().Year()
	if req.BirthYear > currentYear {
		return errors.New("birth year cannot be in the future")
	}
	if req.BirthYear < 1900 {
		return errors.New("birth year is too old")
	}

	// Validate gender
	validGenders := []string{"MALE", "FEMALE", "OTHER"}
	valid := false
	for _, g := range validGenders {
		if req.Gender == g {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("invalid gender value")
	}

	// Validate team
	validTeams := []string{"FRONT_END", "BACK_END", "QA", "ADMIN", "BRSE", "DESIGN", "OTHERS"}
	valid = false
	for _, t := range validTeams {
		if req.Team == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid team value: %s", req.Team)
	}

	return nil
}