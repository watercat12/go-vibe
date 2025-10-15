package dto

import (
	"e-wallet/internal/domain/account"
)

type AccountResponse struct {
	ID              string   `json:"account_id"`
	AccountNumber   string   `json:"account_number"`
	AccountType     string   `json:"account_type"`
	AccountName     string   `json:"account_name"`
	Balance         float64  `json:"balance"`
	InterestRate    *float64 `json:"interest_rate,omitempty"`
	FixedTermMonths *int     `json:"fixed_term_months,omitempty"`
}

type CreateAccountResponse struct {
	Account *AccountResponse `json:"account"`
}

func NewAccountResponse(acc *account.Account) *AccountResponse {
	return &AccountResponse{
		ID:              acc.ID,
		AccountNumber:   acc.AccountNumber,
		AccountType:     acc.AccountType,
		AccountName:     acc.AccountName,
		Balance:         acc.Balance,
		InterestRate:    acc.InterestRate,
		FixedTermMonths: acc.FixedTermMonths,
	}
}

type ListAccountsResponse struct {
	Accounts []*AccountResponse `json:"accounts"`
}

func NewListAccountsResponse(accounts []*account.Account) *ListAccountsResponse {
	accountResponses := make([]*AccountResponse, len(accounts))
	for i, acc := range accounts {
		accountResponses[i] = NewAccountResponse(acc)
	}
	return &ListAccountsResponse{
		Accounts: accountResponses,
	}
}
func NewCreateAccountResponse(acc *account.Account) *CreateAccountResponse {
	return &CreateAccountResponse{
		Account: NewAccountResponse(acc),
	}
}