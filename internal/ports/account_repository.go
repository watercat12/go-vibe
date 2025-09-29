package ports

import (
	"context"
	"e-wallet/internal/domain/account"
)

type AccountRepository interface {
	Create(ctx context.Context, account *account.Account) (*account.Account, error)
	GetByUserID(ctx context.Context, userID string) (*account.Account, error)
	GetByID(ctx context.Context, id string) (*account.Account, error)
}