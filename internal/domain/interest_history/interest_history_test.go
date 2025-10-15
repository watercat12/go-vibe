package interest_history

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInterestHistory(t *testing.T) {
	accountID := "acc-123"
	date := time.Now()
	interestAmount := 100.5

	history := NewInterestHistory(accountID, date, interestAmount)

	assert.NotEmpty(t, history.ID)
	assert.Equal(t, accountID, history.AccountID)
	assert.Equal(t, date, history.Date)
	assert.Equal(t, interestAmount, history.InterestAmount)
}