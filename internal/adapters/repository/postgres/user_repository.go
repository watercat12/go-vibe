package postgres

import (
	"context"
	"errors"

	"e-wallet/internal/domain/user"
	"e-wallet/internal/ports"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *user.User) (*user.User, error) {
	schema := &User{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		PasswordHash:    user.PasswordHash,
		IsEmailVerified: user.IsEmailVerified,
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