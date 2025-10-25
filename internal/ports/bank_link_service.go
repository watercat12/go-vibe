package ports

import (
	"context"

	"e-wallet/internal/domain/bank_link"
)

type BankLinkService interface {
	LinkBankAccount(ctx context.Context, userID, bankCode, accountType string) (*bank_link.BankLink, error)
}