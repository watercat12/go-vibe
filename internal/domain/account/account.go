package account

import (
	"context"
	"time"
)

const (
	PaymentAccountType        = "payment"
	FixedSavingsAccountType   = "savings_fixed"
	FlexibleSavingsAccountType = "savings_flexible"
)

type Account struct {
	ID               string    `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	UserID           string    `json:"user_id" gorm:"not null"`
	AccountType      string    `json:"account_type" gorm:"not null"`
	AccountNumber    string    `json:"account_number" gorm:"unique;not null"`
	AccountName      string    `json:"account_name"`
	Balance          float64   `json:"balance" gorm:"type:numeric(18,2);default:0.00"`
	InterestRate     *float64  `json:"interest_rate" gorm:"type:numeric(5,2)"`
	FixedTermMonths  *int      `json:"fixed_term_months"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateAccountRequest struct {
	UserID      string `json:"user_id"`
	AccountType string `json:"account_type"`
}

type AccountService interface {
	CreatePaymentAccount(ctx context.Context, userID string) (*Account, error)
	CreateFixedSavingsAccount(ctx context.Context, userID string, termMonths int) (*Account, error)
	CreateFlexibleSavingsAccount(ctx context.Context, userID string) (*Account, error)
	CalculateDailyInterest(ctx context.Context) error
}