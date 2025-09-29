package postgres

import (
	"context"
	"errors"

	"e-wallet/internal/domain/user"
	"e-wallet/internal/ports"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ports.ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) Create(ctx context.Context, profile *user.Profile) (*user.Profile, error) {
	schema := &Profile{
		ID:          profile.ID,
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		PhoneNumber: profile.PhoneNumber,
		NationalID:  profile.NationalID,
		BirthYear:   profile.BirthYear,
		Gender:      profile.Gender,
		Team:        profile.Team,
	}

	if err := r.db.WithContext(ctx).Table(ProfilesTableName).Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*user.Profile, error) {
	var schema Profile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *profileRepository) Update(ctx context.Context, profile *user.Profile) (*user.Profile, error) {
	schema := &Profile{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		PhoneNumber: profile.PhoneNumber,
		NationalID:  profile.NationalID,
		BirthYear:   profile.BirthYear,
		Gender:      profile.Gender,
		Team:        profile.Team,
	}

	if err := r.db.WithContext(ctx).Table(ProfilesTableName).Where("user_id = ?", profile.UserID).Updates(schema).Error; err != nil {
		return nil, err
	}

	return r.GetByUserID(ctx, profile.UserID)
}

func (r *profileRepository) Upsert(ctx context.Context, profile *user.Profile) (*user.Profile, error) {
	schema := &Profile{
		ID:          profile.ID,
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		PhoneNumber: profile.PhoneNumber,
		NationalID:  profile.NationalID,
		BirthYear:   profile.BirthYear,
		Gender:      profile.Gender,
		Team:        profile.Team,
	}

	err := r.db.WithContext(ctx).Table(ProfilesTableName).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(schema).Error
	if err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}