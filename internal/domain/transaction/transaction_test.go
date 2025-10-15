package transaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInterestTransaction(t *testing.T) {
	accountID := "acc-123"
	amount := 100.5
	balanceAfter := 1000.0

	tx := NewInterestTransaction(accountID, amount, balanceAfter)

	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, accountID, tx.AccountID)
	assert.Equal(t, TransactionTypeInterest, tx.TransactionType)
	assert.Equal(t, amount, tx.Amount)
	assert.Equal(t, TransactionStatusSuccess, tx.Status)
	assert.Equal(t, balanceAfter, tx.BalanceAfter)
	assert.Nil(t, tx.RelatedAccountID)
}