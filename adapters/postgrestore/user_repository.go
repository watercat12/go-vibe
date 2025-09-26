package postgrestore

import (
	"context"
	"errors"
	"fmt"

	"e-wallet/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	schema := &User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	if err := r.db.WithContext(ctx).Table(UsersTableName).Create(schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var schema User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var schema User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, profile *user.UserProfile) error {
	schema := UserProfile{
		Name:      profile.Name,
		Email:     profile.Email,
		Avatar:    profile.Avatar,
		Phone:     profile.Phone,
		IDNumber:  profile.IDNumber,
		BirthYear: profile.BirthYear,
		Gender:    profile.Gender,
		Team:      profile.Team,
	}

	// Upsert: if exists update, else create
	if err := r.db.WithContext(ctx).Table(UserProfilesTableName).
	Where("user_id = ?", profile.UserID).
	Assign(schema).
	FirstOrCreate(&schema).Error; err != nil {
		return err
	}

	fmt.Printf("schema: %v\n", schema)

	return nil
}