package interest_history

import (
	"time"
)

type InterestHistory struct {
	ID             string    `json:"id" gorm:"type:uuid;default:gen_random_uuid()"`
	AccountID      string    `json:"account_id" gorm:"not null"`
	Date           time.Time `json:"date" gorm:"type:date;not null"`
	InterestAmount float64   `json:"interest_amount" gorm:"type:numeric(18,2);not null"`
	CreatedAt      time.Time `json:"created_at"`
}