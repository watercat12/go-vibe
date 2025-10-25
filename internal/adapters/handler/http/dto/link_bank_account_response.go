package dto

import "e-wallet/internal/domain/bank_link"

type LinkBankAccountResponse struct {
	ID          string `json:"id"`
	BankCode    string `json:"bank_code"`
	AccountType string `json:"account_type"`
	Status      string `json:"status"`
}

func NewLinkBankAccountResponse(bankLink *bank_link.BankLink) *LinkBankAccountResponse {
	return &LinkBankAccountResponse{
		ID:          bankLink.ID,
		BankCode:    bankLink.BankCode,
		AccountType: bankLink.AccountType,
		Status:      bankLink.Status,
	}
}