package transaction

import (
	"time"
)

type Transaction struct {
	ID               string    `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	AccountID        string    `json:"account_id" gorm:"not null"`
	TransactionType  string    `json:"transaction_type" gorm:"not null"`
	Amount           float64   `json:"amount" gorm:"type:numeric(18,2);not null"`
	Status           string    `json:"status" gorm:"not null"`
	BalanceAfter     float64   `json:"balance_after" gorm:"type:numeric(18,2);not null"`
	RelatedAccountID *string   `json:"related_account_id"`
	CreatedAt        time.Time `json:"created_at"`
}

const (
	TransactionTypeInterest = "interest"
	TransactionStatusSuccess = "success"
)