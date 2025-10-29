package ports

import (
	"context"
	"e-wallet/internal/domain/account"
)

type AccountService interface {
	CreatePaymentAccount(ctx context.Context, userID string) (*account.Account, error)
	CreateFixedSavingsAccount(ctx context.Context, userID string, req *account.CreateFixedSavingsAccountRequest) (*account.Account, error)
	CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*account.Account, error)
	ListAccounts(ctx context.Context, userID string) (*account.ListAccountsResponse, error)
}