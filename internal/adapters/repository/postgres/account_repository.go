package postgres

import (
	"context"
	"errors"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/ports"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) ports.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *account.Account) (*account.Account, error) {
	schema := Account{
		ID:              account.ID,
		UserID:          account.UserID,
		AccountType:     account.AccountType,
		AccountNumber:   account.AccountNumber,
		AccountName:     account.AccountName,
		Balance:         account.Balance,
		InterestRate:    account.InterestRate,
		FixedTermMonths: account.FixedTermMonths,
	}

	if err := r.db.WithContext(ctx).Table(AccountsTableName).Create(&schema).Error; err != nil {
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID string) (*account.Account, error) {
	var schema Account
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*account.Account, error) {
	var schema Account
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schema).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return schema.ToDomain(), nil
}