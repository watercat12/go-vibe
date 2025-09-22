package postgrestore

import (
	"context"
	"errors"

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
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
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

func (r *userRepository) GetByID(ctx context.Context, id int) (*user.User, error) {
	var schema User
	if err := r.db.WithContext(ctx).First(&schema, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}