package bank_link

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBankLink(t *testing.T) {
	userID := "user-123"
	bankCode := "BANK001"
	accountType := "SAVINGS"
	accessToken := "access-token"
	refreshToken := "refresh-token"
	expiresIn := 3600

	bankLink := NewBankLink(userID, bankCode, accountType, accessToken, refreshToken, expiresIn)

	assert.NotEmpty(t, bankLink.ID)
	assert.Equal(t, userID, bankLink.UserID)
	assert.Equal(t, bankCode, bankLink.BankCode)
	assert.Equal(t, accountType, bankLink.AccountType)
	assert.Equal(t, accessToken, bankLink.AccessToken)
	assert.Equal(t, refreshToken, bankLink.RefreshToken)
	assert.Equal(t, expiresIn, bankLink.ExpiresIn)
	assert.Equal(t, "ACTIVE", bankLink.Status)
	assert.NotZero(t, bankLink.CreatedAt)
	assert.WithinDuration(t, time.Now(), bankLink.CreatedAt, time.Second)
}

func TestNewBankLinkResponse(t *testing.T) {
	bankLink := &BankLink{
		ID:          "link-123",
		UserID:      "user-123",
		BankCode:    "BANK001",
		AccountType: "SAVINGS",
		AccessToken: "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
		Status:       "ACTIVE",
		CreatedAt:    time.Now(),
	}

	response := NewBankLinkResponse(bankLink)

	assert.Equal(t, bankLink.ID, response.ID)
	assert.Equal(t, bankLink.BankCode, response.BankCode)
	assert.Equal(t, bankLink.AccountType, response.AccountType)
	assert.Equal(t, bankLink.Status, response.Status)
}