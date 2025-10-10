package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/interest_history"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
)

func TestInterestHistoryRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	accountRepo := NewAccountRepository(db)
	repo := NewInterestHistoryRepository(db)

	tests := []struct {
		name            string
		setupAccount    func() string
		interestHistory func(accountID string) *interest_history.InterestHistory
		expectError     bool
	}{
		{
			name: "success - create interest history",
			setupAccount: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser1",
					Email:        "test1@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				account := &account.Account{
					ID:            pkg.NewUUIDV7(),
					UserID:        user.ID,
					AccountType:   account.PaymentAccountType,
					AccountNumber: "PAY123456789",
					AccountName:   "Payment Account",
					Balance:       1000.0,
				}
				_, err = accountRepo.Create(context.Background(), account)
				require.NoError(t, err)
				return account.ID
			},
			interestHistory: func(accountID string) *interest_history.InterestHistory {
				return &interest_history.InterestHistory{
					ID:             pkg.NewUUIDV7(),
					AccountID:      accountID,
					Date:           time.Now(),
					InterestAmount: 100.50,
				}
			},
			expectError: false,
		},
		{
			name: "success - create interest history with zero amount",
			setupAccount: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser2",
					Email:        "test2@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				account := &account.Account{
					ID:            pkg.NewUUIDV7(),
					UserID:        user.ID,
					AccountType:   account.FlexibleSavingsAccountType,
					AccountNumber: "SAV123456789",
					AccountName:   "Flexible Savings",
					Balance:       5000.0,
				}
				_, err = accountRepo.Create(context.Background(), account)
				require.NoError(t, err)
				return account.ID
			},
			interestHistory: func(accountID string) *interest_history.InterestHistory {
				return &interest_history.InterestHistory{
					ID:             pkg.NewUUIDV7(),
					AccountID:      accountID,
					Date:           time.Now().Add(-24 * time.Hour),
					InterestAmount: 0.0,
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountID := tt.setupAccount()
			testInterestHistory := tt.interestHistory(accountID)
			result, err := repo.Create(context.Background(), testInterestHistory)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testInterestHistory.ID, result.ID)
				assert.Equal(t, testInterestHistory.AccountID, result.AccountID)
				assert.Equal(t, testInterestHistory.Date, result.Date)
				assert.Equal(t, testInterestHistory.InterestAmount, result.InterestAmount)
				assert.NotZero(t, result.CreatedAt)
			}
		})
	}
}

func TestInterestHistoryRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewInterestHistoryRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	interestHistory := &interest_history.InterestHistory{
		ID:             pkg.NewUUIDV7(),
		AccountID:      pkg.NewUUIDV7(),
		Date:           time.Now(),
		InterestAmount: 50.25,
	}

	result, err := repo.Create(context.Background(), interestHistory)

	assert.Error(t, err)
	assert.Nil(t, result)
}