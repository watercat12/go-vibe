package dto

import (
	"e-wallet/internal/domain/account"
)

type AccountResponse struct {
	ID            string  `json:"account_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}

type CreateAccountResponse struct {
	Account *AccountResponse `json:"account"`
}

func NewAccountResponse(acc *account.Account) *AccountResponse {
	return &AccountResponse{
		ID:            acc.ID,
		AccountNumber: acc.AccountNumber,
		Balance:       acc.Balance,
	}
}

func NewCreateAccountResponse(acc *account.Account) *CreateAccountResponse {
	return &CreateAccountResponse{
		Account: NewAccountResponse(acc),
	}
}