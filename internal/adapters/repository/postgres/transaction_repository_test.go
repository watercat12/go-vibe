package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/transaction"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
)

func TestTransactionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	accountRepo := NewAccountRepository(db)
	transactionRepo := NewTransactionRepository(db)

	// Setup test user and account
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	testAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        testUser.ID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY123456789",
		AccountName:   "Payment Account",
		Balance:       1000.0,
	}
	_, err = accountRepo.Create(context.Background(), testAccount)
	require.NoError(t, err)

	// Setup another account for transfer
	relatedAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        testUser.ID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY987654321",
		AccountName:   "Related Account",
		Balance:       500.0,
	}
	_, err = accountRepo.Create(context.Background(), relatedAccount)
	require.NoError(t, err)

	tests := []struct {
		name        string
		transaction func() *transaction.Transaction
		expectError bool
	}{
		{
			name: "success - create interest transaction",
			transaction: func() *transaction.Transaction {
				return &transaction.Transaction{
					ID:              pkg.NewUUIDV7(),
					AccountID:       testAccount.ID,
					TransactionType: transaction.TransactionTypeInterest,
					Amount:          50.0,
					Status:          transaction.TransactionStatusSuccess,
					BalanceAfter:    1050.0,
				}
			},
			expectError: false,
		},
		{
			name: "success - create deposit transaction",
			transaction: func() *transaction.Transaction {
				return &transaction.Transaction{
					ID:              pkg.NewUUIDV7(),
					AccountID:       testAccount.ID,
					TransactionType: "deposit",
					Amount:          200.0,
					Status:          transaction.TransactionStatusSuccess,
					BalanceAfter:    1200.0,
				}
			},
			expectError: false,
		},
		{
			name: "success - create withdraw transaction",
			transaction: func() *transaction.Transaction {
				return &transaction.Transaction{
					ID:              pkg.NewUUIDV7(),
					AccountID:       testAccount.ID,
					TransactionType: "withdraw",
					Amount:          100.0,
					Status:          transaction.TransactionStatusSuccess,
					BalanceAfter:    900.0,
				}
			},
			expectError: false,
		},
		{
			name: "success - create transfer transaction",
			transaction: func() *transaction.Transaction {
				relatedID := relatedAccount.ID
				return &transaction.Transaction{
					ID:               pkg.NewUUIDV7(),
					AccountID:        testAccount.ID,
					TransactionType:  "transfer",
					Amount:           150.0,
					Status:           transaction.TransactionStatusSuccess,
					BalanceAfter:     850.0,
					RelatedAccountID: &relatedID,
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTx := tt.transaction()
			result, err := transactionRepo.Create(context.Background(), testTx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testTx.ID, result.ID)
				assert.Equal(t, testTx.AccountID, result.AccountID)
				assert.Equal(t, testTx.TransactionType, result.TransactionType)
				assert.Equal(t, testTx.Amount, result.Amount)
				assert.Equal(t, testTx.Status, result.Status)
				assert.Equal(t, testTx.BalanceAfter, result.BalanceAfter)
				if testTx.RelatedAccountID != nil {
					assert.Equal(t, *testTx.RelatedAccountID, *result.RelatedAccountID)
				} else {
					assert.Nil(t, result.RelatedAccountID)
				}
				assert.NotZero(t, result.CreatedAt)
			}
		})
	}
}

func TestTransactionRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	tx := &transaction.Transaction{
		ID:              pkg.NewUUIDV7(),
		AccountID:       pkg.NewUUIDV7(),
		TransactionType: transaction.TransactionTypeInterest,
		Amount:          50.0,
		Status:          transaction.TransactionStatusSuccess,
		BalanceAfter:    1050.0,
	}

	result, err := repo.Create(context.Background(), tx)

	assert.Error(t, err)
	assert.Nil(t, result)
}