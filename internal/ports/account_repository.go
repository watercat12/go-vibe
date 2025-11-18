package ports

import (
	"context"
	"e-wallet/internal/domain/account"
	"time"
)

type AccountRepository interface {
	CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error)
	CreateFixedSavingsAccount(ctx context.Context, userID string, req *account.CreateFixedSavingsAccountRequest) (*account.Account, error)
	CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]*account.Account, error)
	GetAccountByID(ctx context.Context, accountID string) (*account.Account, error)
	CountPaymentAccountsByUserID(ctx context.Context, userID string) (int64, error)
	CountSavingsAccountsByUserID(ctx context.Context, userID string) (int64, error)
	UpdateAccountBalance(ctx context.Context, accountID string, newBalance float64) error
}

type SavingsAccountDetailRepository interface {
	CreateSavingsAccountDetail(ctx context.Context, detail *account.SavingsAccountDetail) error
	GetSavingsAccountDetailByAccountID(ctx context.Context, accountID string) (*account.SavingsAccountDetail, error)
	UpdateLastInterestCalcDate(ctx context.Context, accountID string, date *time.Time) error
}