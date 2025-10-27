package postgres

import (
	"context"
	"errors"

	"e-wallet/internal/domain/profile"
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

func (r *profileRepository) GetByUserID(ctx context.Context, userID string) (*profile.Profile, error) {
	var schema UserProfile
	if err := r.db.WithContext(ctx).Table(UserProfilesTableName).Where("user_id = ?", userID).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *profileRepository) Upsert(ctx context.Context, profile *profile.Profile) (*profile.Profile, error) {
	schema := &UserProfile{
		UserID:      profile.UserID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL,
		PhoneNumber: profile.PhoneNumber,
		NationalID:  profile.NationalID,
		BirthYear:   profile.BirthYear,
		Gender:      profile.Gender,
		Team:        profile.Team,
	}

	if err := r.db.WithContext(ctx).Table(UserProfilesTableName).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"display_name", "avatar_url", "phone_number", "national_id", "birth_year", "gender", "team", "updated_at"}),
		}).
		Clauses(clause.Returning{}).
		Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *profileRepository) CheckNationalIDExists(ctx context.Context, nationalID string, excludeUserID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Table(UserProfilesTableName).Where("national_id = ?", nationalID)
	if excludeUserID != "" {
		query = query.Where("user_id != ?", excludeUserID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
