package ports

import (
	"context"
	"e-wallet/internal/domain/transaction"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error)
}