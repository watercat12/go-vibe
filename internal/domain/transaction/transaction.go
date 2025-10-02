package transaction

import (
	"time"

	"e-wallet/pkg"
)

type Transaction struct {
	ID               string
	AccountID        string
	TransactionType  string
	Amount           float64
	Status           string
	BalanceAfter     float64
	RelatedAccountID *string
	CreatedAt        time.Time
}

const (
	TransactionTypeInterest = "interest"
	TransactionStatusSuccess = "success"
)

func NewInterestTransaction(accountID string, amount float64, balanceAfter float64) *Transaction {
	return &Transaction{
		ID:              pkg.NewUUIDV7(),
		AccountID:       accountID,
		TransactionType: TransactionTypeInterest,
		Amount:          amount,
		Status:          TransactionStatusSuccess,
		BalanceAfter:    balanceAfter,
	}
}