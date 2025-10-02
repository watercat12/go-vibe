package transaction

import (
	"time"
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