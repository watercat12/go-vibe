package account

import (
	"time"
)

const (
	PaymentAccountType        = "payment"
	FixedSavingsAccountType   = "savings_fixed"
	FlexibleSavingsAccountType = "savings_flexible"
)

type Account struct {
	ID               string
	UserID           string
	AccountType      string
	AccountNumber    string
	AccountName      string
	Balance          float64
	InterestRate     *float64
	FixedTermMonths  *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CreateAccountRequest struct {
	UserID      string
	AccountType string
}