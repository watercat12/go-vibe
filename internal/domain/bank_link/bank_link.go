package bank_link

import (
	"time"

	"github.com/google/uuid"
)

type BankLink struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	BankCode     string    `json:"bank_code"`
	AccountType  string    `json:"account_type"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateBankLinkRequest struct {
	UserID      string `json:"user_id"`
	BankCode    string `json:"bank_code"`
	AccountType string `json:"account_type"`
}

type BankLinkResponse struct {
	ID          string `json:"id"`
	BankCode    string `json:"bank_code"`
	AccountType string `json:"account_type"`
	Status      string `json:"status"`
}

func NewBankLink(userID, bankCode, accountType, accessToken, refreshToken string, expiresIn int) *BankLink {
	return &BankLink{
		ID:           uuid.New().String(),
		UserID:       userID,
		BankCode:     bankCode,
		AccountType:  accountType,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
	}
}

func NewBankLinkResponse(bankLink *BankLink) *BankLinkResponse {
	return &BankLinkResponse{
		ID:          bankLink.ID,
		BankCode:    bankLink.BankCode,
		AccountType: bankLink.AccountType,
		Status:      bankLink.Status,
	}
}