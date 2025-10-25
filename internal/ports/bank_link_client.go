package ports

import (
	"context"

	"e-wallet/internal/domain/bank_link"
)

type BankLinkClient interface {
	LinkAccount(ctx context.Context, bankCode, accountType string) (*bank_link.BankLink, error)
}