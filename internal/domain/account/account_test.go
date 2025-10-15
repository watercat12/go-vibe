package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateTermMonths(t *testing.T) {
	tests := []struct {
		name         string
		termMonths   int
		expectedRate float64
		expectError  bool
	}{
		{
			name:         "success - valid term 3 months",
			termMonths:   3,
			expectedRate: 1.8,
			expectError:  false,
		},
		{
			name:         "success - valid term 12 months",
			termMonths:   12,
			expectedRate: 7.2,
			expectError:  false,
		},
		{
			name:        "error - invalid term",
			termMonths: 2,
			expectError: true,
		},
		{
			name:        "error - invalid term 0",
			termMonths: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, err := ValidateTermMonths(tt.termMonths)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0.0, rate)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRate, rate)
			}
		})
	}
}

func TestAccount_CalculateDailyInterest(t *testing.T) {
	tests := []struct {
		name           string
		balance        float64
		createdAt      time.Time
		expectedDaily  float64
	}{
		{
			name:          "success - under 10M balance, within 30 days",
			balance:       5000000,
			createdAt:     time.Now().Add(-10 * 24 * time.Hour),
			expectedDaily: 109.589,
		},
		{
			name:          "success - under 10M balance, over 30 days",
			balance:       5000000,
			createdAt:     time.Now().Add(-40 * 24 * time.Hour),
			expectedDaily: 41.096,
		},
		{
			name:          "success - 10-50M balance",
			balance:       20000000,
			createdAt:     time.Now().Add(-40 * 24 * time.Hour),
			expectedDaily: 219.178,
		},
		{
			name:          "success - over 50M balance",
			balance:       60000000,
			createdAt:     time.Now().Add(-40 * 24 * time.Hour),
			expectedDaily: 821.918,
		},
		{
			name:          "success - zero balance",
			balance:       0,
			createdAt:     time.Now().Add(-40 * 24 * time.Hour),
			expectedDaily: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &Account{
				Balance:   tt.balance,
				CreatedAt: tt.createdAt,
			}
			daily := acc.CalculateDailyInterest()
			assert.Equal(t, tt.expectedDaily, daily)
		})
	}
}

func TestGenerateAccountNumber(t *testing.T) {
	tests := []struct {
		name        string
		accountType string
		expectedPrefix string
	}{
		{
			name:           "success - payment account",
			accountType:    PaymentAccountType,
			expectedPrefix: "PAY",
		},
		{
			name:           "success - fixed savings account",
			accountType:    FixedSavingsAccountType,
			expectedPrefix: "SAV",
		},
		{
			name:           "success - flexible savings account",
			accountType:    FlexibleSavingsAccountType,
			expectedPrefix: "SAV",
		},
		{
			name:           "success - unknown type defaults to PAY",
			accountType:    "unknown",
			expectedPrefix: "PAY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			number := GenerateAccountNumber(tt.accountType)
			assert.Contains(t, number, tt.expectedPrefix)
			assert.Greater(t, len(number), len(tt.expectedPrefix))
		})
	}
}

func TestNewPaymentAccount(t *testing.T) {
	userID := "user-123"
	acc := NewPaymentAccount(userID)

	assert.NotEmpty(t, acc.ID)
	assert.Equal(t, userID, acc.UserID)
	assert.Equal(t, PaymentAccountType, acc.AccountType)
	assert.Contains(t, acc.AccountNumber, "PAY")
	assert.Equal(t, 0.0, acc.Balance)
	assert.Nil(t, acc.InterestRate)
	assert.Nil(t, acc.FixedTermMonths)
}

func TestNewFlexibleSavingsAccount(t *testing.T) {
	userID := "user-123"
	acc := NewFlexibleSavingsAccount(userID)

	assert.NotEmpty(t, acc.ID)
	assert.Equal(t, userID, acc.UserID)
	assert.Equal(t, FlexibleSavingsAccountType, acc.AccountType)
	assert.Contains(t, acc.AccountNumber, "SAV")
	assert.Equal(t, 0.0, acc.Balance)
	assert.Nil(t, acc.InterestRate)
	assert.Nil(t, acc.FixedTermMonths)
}

func TestNewFixedSavingsAccount(t *testing.T) {
	userID := "user-123"
	termMonths := 3
	interestRate := 1.8
	acc := NewFixedSavingsAccount(userID, termMonths, interestRate)

	assert.NotEmpty(t, acc.ID)
	assert.Equal(t, userID, acc.UserID)
	assert.Equal(t, FixedSavingsAccountType, acc.AccountType)
	assert.Contains(t, acc.AccountNumber, "SAV")
	assert.Equal(t, 0.0, acc.Balance)
	assert.NotNil(t, acc.InterestRate)
	assert.Equal(t, interestRate, *acc.InterestRate)
	assert.NotNil(t, acc.FixedTermMonths)
	assert.Equal(t, termMonths, *acc.FixedTermMonths)
}