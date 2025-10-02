package interest_history

import (
	"time"

	"e-wallet/pkg"
)

type InterestHistory struct {
	ID             string
	AccountID      string
	Date           time.Time
	InterestAmount float64
	CreatedAt      time.Time
}

func NewInterestHistory(accountID string, date time.Time, interestAmount float64) *InterestHistory {
	return &InterestHistory{
		ID:             pkg.NewUUIDV7(),
		AccountID:      accountID,
		Date:           date,
		InterestAmount: interestAmount,
	}
}