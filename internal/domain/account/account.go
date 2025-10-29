package account

import (
	"time"
)

type Account struct {
	ID            string
	UserID        string
	AccountNumber string
	AccountType   string
	Balance       float64
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreatePaymentAccountRequest struct{}

type CreateFixedSavingsAccountRequest struct {
	TermCode string `json:"term_code" validate:"required,oneof=1 3 6 8 12"`
}

type CreateFlexibleSavingsAccountRequest struct{}

type AccountResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SavingsAccountDetail struct {
	AccountID             string
	IsFixedTerm           bool
	TermMonths            *int
	AnnualInterestRate    float64
	StartDate             time.Time
	MaturityDate          *time.Time
	LastInterestCalcDate  *time.Time
}

type SavingsAccountDetailResponse struct {
	AccountID             string     `json:"account_id"`
	IsFixedTerm           bool       `json:"is_fixed_term"`
	TermMonths            *int       `json:"term_months,omitempty"`
	AnnualInterestRate    float64    `json:"annual_interest_rate"`
	StartDate             time.Time  `json:"start_date"`
	MaturityDate          *time.Time `json:"maturity_date,omitempty"`
	LastInterestCalcDate  *time.Time `json:"last_interest_calc_date,omitempty"`
}

type ListAccountsResponse struct {
	Accounts []AccountWithDetailsResponse `json:"accounts"`
}

type AccountWithDetailsResponse struct {
	AccountResponse
	SavingsDetail *SavingsAccountDetailResponse `json:"savings_detail,omitempty"`
}