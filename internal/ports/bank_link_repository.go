package ports

import (
	"context"

	"e-wallet/internal/domain/bank_link"
)

type BankLinkRepository interface {
	Create(ctx context.Context, bankLink *bank_link.BankLink) (*bank_link.BankLink, error)
	GetByUserID(ctx context.Context, userID string) ([]*bank_link.BankLink, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
}