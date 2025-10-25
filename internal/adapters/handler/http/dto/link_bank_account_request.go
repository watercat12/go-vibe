package dto

type LinkBankAccountRequest struct {
	BankCode    string `json:"bank_code" validate:"required"`
	AccountType string `json:"account_type" validate:"required"`
}