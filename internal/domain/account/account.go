package account

import (
	"errors"
	"fmt"
	"time"

	"e-wallet/pkg"
)

const (
	PaymentAccountType        = "payment"
	FixedSavingsAccountType   = "savings_fixed"
	FlexibleSavingsAccountType = "savings_flexible"
)

var ValidTerms = map[int]float64{
	1:  0.6,
	3:  1.8,
	6:  3.6,
	8:  4.8,
	12: 7.2,
}

func ValidateTermMonths(termMonths int) (float64, error) {
	interestRate, ok := ValidTerms[termMonths]
	if !ok {
		return 0, errors.New("invalid term months")
	}
	return interestRate, nil
}

type Account struct {
	ID               string
	UserID           string
	AccountType      string
	AccountNumber    string
	AccountName      string
	Balance          float64
	InterestRate     *float64
	FixedTermMonths  *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (a *Account) CalculateDailyInterest() float64 {
	var annualRate float64
	switch {
	case time.Since(a.CreatedAt) < 30*24*time.Hour:
		annualRate = 0.008
	case a.Balance < 10000000: // under 10M
		annualRate = 0.003
	case a.Balance < 50000000: // 10-50M
		annualRate = 0.004
	default: // over 50M
		annualRate = 0.005
	}

	dailyRate := annualRate / 365
	return a.Balance * dailyRate
}

func GenerateAccountNumber(accountType string) string {
	prefix := "PAY"
	if accountType == FixedSavingsAccountType || accountType == FlexibleSavingsAccountType {
		prefix = "SAV"
	}
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
}

func NewPaymentAccount(userID string) *Account {
	return &Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountType:   PaymentAccountType,
		AccountNumber: GenerateAccountNumber(PaymentAccountType),
		Balance:       0.0,
	}
}

func NewFlexibleSavingsAccount(userID string) *Account {
	return &Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        userID,
		AccountType:   FlexibleSavingsAccountType,
		AccountNumber: GenerateAccountNumber(FlexibleSavingsAccountType),
		Balance:       0.0,
	}
}

func NewFixedSavingsAccount(userID string, termMonths int, interestRate float64) *Account {
	return &Account{
		ID:              pkg.NewUUIDV7(),
		UserID:          userID,
		AccountType:     FixedSavingsAccountType,
		AccountNumber:   GenerateAccountNumber(FixedSavingsAccountType),
		Balance:         0.0,
		InterestRate:    &interestRate,
		FixedTermMonths: &termMonths,
	}
}

type CreateAccountRequest struct {
	UserID      string
	AccountType string
}