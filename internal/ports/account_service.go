package ports

import (
	"context"

	"e-wallet/internal/domain/account"
)

type AccountService interface {
	CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error)
	CreateFixedSavingsAccount(ctx context.Context, userID string, termMonths int) (*account.Account, error)
	CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error)
	GetAccountsByUserID(ctx context.Context, userID string) ([]*account.Account, error)
	CalculateDailyInterest(ctx context.Context) error
}