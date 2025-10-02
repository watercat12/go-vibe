package interest_history

import (
	"time"
)

type InterestHistory struct {
	ID             string
	AccountID      string
	Date           time.Time
	InterestAmount float64
	CreatedAt      time.Time
}