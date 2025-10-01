package ports

import (
	"context"
	"e-wallet/internal/domain/account"
)

type AccountRepository interface {
	Create(ctx context.Context, account *account.Account) (*account.Account, error)
	GetByUserID(ctx context.Context, userID string) (*account.Account, error)
	GetByID(ctx context.Context, id string) (*account.Account, error)
	CountSavingsAccounts(ctx context.Context, userID string) (int64, error)
	GetFlexibleSavingsAccounts(ctx context.Context) ([]*account.Account, error)
	UpdateBalance(ctx context.Context, id string, balance float64) error
}