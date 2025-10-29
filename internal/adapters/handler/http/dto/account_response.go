package dto

import (
	"e-wallet/internal/domain/account"
	"time"
)

type AccountResponse struct {
	ID            string    `json:"id" example:"acc-123"`
	UserID        string    `json:"user_id" example:"user-123"`
	AccountNumber string    `json:"account_number" example:"1234567890"`
	AccountType   string    `json:"account_type" example:"payment"`
	Balance       float64   `json:"balance" example:"1000.50"`
	Status        string    `json:"status" example:"active"`
	CreatedAt     time.Time `json:"created_at" example:"2023-10-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2023-10-01T00:00:00Z"`
}

func NewAccountResponse(acc *account.Account) *AccountResponse {
	return &AccountResponse{
		ID:            acc.ID,
		UserID:        acc.UserID,
		AccountNumber: acc.AccountNumber,
		AccountType:   acc.AccountType,
		Balance:       acc.Balance,
		Status:        acc.Status,
		CreatedAt:     acc.CreatedAt,
		UpdatedAt:     acc.UpdatedAt,
	}
}

type SavingsAccountDetailResponse struct {
	AccountID             string     `json:"account_id" example:"acc-123"`
	IsFixedTerm           bool       `json:"is_fixed_term" example:"true"`
	TermMonths            *int       `json:"term_months,omitempty" example:"12"`
	AnnualInterestRate    float64    `json:"annual_interest_rate" example:"5.5"`
	StartDate             time.Time  `json:"start_date" example:"2023-10-01T00:00:00Z"`
	MaturityDate          *time.Time `json:"maturity_date,omitempty" example:"2024-10-01T00:00:00Z"`
	LastInterestCalcDate  *time.Time `json:"last_interest_calc_date,omitempty" example:"2023-10-01T00:00:00Z"`
}

type AccountWithDetailsResponse struct {
	AccountResponse
	SavingsDetail *SavingsAccountDetailResponse `json:"savings_detail,omitempty"`
}

type ListAccountsResponse struct {
	Accounts []AccountWithDetailsResponse `json:"accounts"`
}