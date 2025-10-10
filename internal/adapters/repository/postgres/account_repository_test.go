package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
)

func TestAccountRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	accountRepo := NewAccountRepository(db)

	tests := []struct {
		name        string
		setupUser   func() string
		account     func(userID string) *account.Account
		expectError bool
	}{
		{
			name: "success - create payment account",
			setupUser: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser1",
					Email:        "test1@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				return user.ID
			},
			account: func(userID string) *account.Account {
				return &account.Account{
					ID:            pkg.NewUUIDV7(),
					UserID:        userID,
					AccountType:   account.PaymentAccountType,
					AccountNumber: "PAY123456789",
					AccountName:   "Payment Account",
					Balance:       1000.0,
				}
			},
			expectError: false,
		},
		{
			name: "success - create flexible savings account",
			setupUser: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser2",
					Email:        "test2@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				return user.ID
			},
			account: func(userID string) *account.Account {
				return &account.Account{
					ID:            pkg.NewUUIDV7(),
					UserID:        userID,
					AccountType:   account.FlexibleSavingsAccountType,
					AccountNumber: "SAV123456789",
					AccountName:   "Flexible Savings",
					Balance:       5000.0,
				}
			},
			expectError: false,
		},
		{
			name: "success - create fixed savings account",
			setupUser: func() string {
				user := &user.User{
					ID:           pkg.NewUUIDV7(),
					Username:     "testuser3",
					Email:        "test3@example.com",
					PasswordHash: "hashedpassword",
				}
				_, err := userRepo.Create(context.Background(), user)
				require.NoError(t, err)
				return user.ID
			},
			account: func(userID string) *account.Account {
				return &account.Account{
					ID:              pkg.NewUUIDV7(),
					UserID:          userID,
					AccountType:     account.FixedSavingsAccountType,
					AccountNumber:   "SAV987654321",
					AccountName:     "Fixed Savings",
					Balance:         10000.0,
					InterestRate:    func() *float64 { r := 5.0; return &r }(),
					FixedTermMonths: func() *int { m := 12; return &m }(),
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := tt.setupUser()
			testAccount := tt.account(userID)
			result, err := accountRepo.Create(context.Background(), testAccount)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testAccount.ID, result.ID)
				assert.Equal(t, testAccount.UserID, result.UserID)
				assert.Equal(t, testAccount.AccountType, result.AccountType)
				assert.Equal(t, testAccount.AccountNumber, result.AccountNumber)
				assert.Equal(t, testAccount.AccountName, result.AccountName)
				assert.Equal(t, testAccount.Balance, result.Balance)
				if testAccount.InterestRate != nil {
					assert.Equal(t, *testAccount.InterestRate, *result.InterestRate)
				}
				if testAccount.FixedTermMonths != nil {
					assert.Equal(t, *testAccount.FixedTermMonths, *result.FixedTermMonths)
				}
				assert.NotZero(t, result.CreatedAt)
				assert.NotZero(t, result.UpdatedAt)
			}
		})
	}
}

func TestAccountRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	account := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        pkg.NewUUIDV7(),
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY123456789",
		AccountName:   "Payment Account",
		Balance:       1000.0,
	}

	result, err := repo.Create(context.Background(), account)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAccountRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewAccountRepository(db)

	// Setup test data
	user := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	testAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        user.ID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY123456789",
		AccountName:   "Payment Account",
		Balance:       1000.0,
	}
	_, err = repo.Create(context.Background(), testAccount)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing account",
			userID:      testAccount.UserID,
			expectError: false,
		},
		{
			name:        "error - account not found",
			userID:      pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByUserID(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testAccount.ID, result.ID)
				assert.Equal(t, testAccount.UserID, result.UserID)
				assert.Equal(t, testAccount.AccountType, result.AccountType)
				assert.Equal(t, testAccount.AccountNumber, result.AccountNumber)
				assert.Equal(t, testAccount.AccountName, result.AccountName)
				assert.Equal(t, testAccount.Balance, result.Balance)
			}
		})
	}
}

func TestAccountRepository_GetByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByUserID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrAccountNotFound, err)
}

func TestAccountRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewAccountRepository(db)

	// Setup test data
	user := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	testAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        user.ID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY123456789",
		AccountName:   "Payment Account",
		Balance:       1000.0,
	}
	_, err = repo.Create(context.Background(), testAccount)
	require.NoError(t, err)

	tests := []struct {
		name        string
		id          string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing account",
			id:          testAccount.ID,
			expectError: false,
		},
		{
			name:        "error - account not found",
			id:          pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByID(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testAccount.ID, result.ID)
				assert.Equal(t, testAccount.UserID, result.UserID)
				assert.Equal(t, testAccount.AccountType, result.AccountType)
				assert.Equal(t, testAccount.AccountNumber, result.AccountNumber)
				assert.Equal(t, testAccount.AccountName, result.AccountName)
				assert.Equal(t, testAccount.Balance, result.Balance)
			}
		})
	}
}

func TestAccountRepository_GetByID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrAccountNotFound, err)
}

func TestAccountRepository_CountSavingsAccounts(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewAccountRepository(db)

	user := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	// Create test accounts
	accounts := []*account.Account{
		{
			ID:            pkg.NewUUIDV7(),
			UserID:        user.ID,
			AccountType:   account.FlexibleSavingsAccountType,
			AccountNumber: "SAV123456789",
			AccountName:   "Flexible Savings 1",
			Balance:       5000.0,
		},
		{
			ID:            pkg.NewUUIDV7(),
			UserID:        user.ID,
			AccountType:   account.FixedSavingsAccountType,
			AccountNumber: "SAV987654321",
			AccountName:   "Fixed Savings 1",
			Balance:       10000.0,
			InterestRate:  func() *float64 { r := 5.0; return &r }(),
			FixedTermMonths: func() *int { m := 12; return &m }(),
		},
		{
			ID:            pkg.NewUUIDV7(),
			UserID:        user.ID,
			AccountType:   account.PaymentAccountType,
			AccountNumber: "PAY111111111",
			AccountName:   "Payment Account",
			Balance:       1000.0,
		},
	}

	for _, acc := range accounts {
		_, err := repo.Create(context.Background(), acc)
		require.NoError(t, err)
	}

	count, err := repo.CountSavingsAccounts(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestAccountRepository_CountSavingsAccounts_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	count, err := repo.CountSavingsAccounts(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAccountRepository_GetFlexibleSavingsAccounts(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewAccountRepository(db)

	user1 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser1",
		Email:        "test1@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), user1)
	require.NoError(t, err)

	user2 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser2",
		Email:        "test2@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err = userRepo.Create(context.Background(), user2)
	require.NoError(t, err)

	// Create test accounts
	flexibleAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        user1.ID,
		AccountType:   account.FlexibleSavingsAccountType,
		AccountNumber: "SAV123456789",
		AccountName:   "Flexible Savings",
		Balance:       5000.0,
	}
	_, err = repo.Create(context.Background(), flexibleAccount)
	require.NoError(t, err)

	// Create a fixed savings account (should not be returned)
	fixedAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        user2.ID,
		AccountType:   account.FixedSavingsAccountType,
		AccountNumber: "SAV987654321",
		AccountName:   "Fixed Savings",
		Balance:       10000.0,
		InterestRate:  func() *float64 { r := 5.0; return &r }(),
		FixedTermMonths: func() *int { m := 12; return &m }(),
	}
	_, err = repo.Create(context.Background(), fixedAccount)
	require.NoError(t, err)

	accounts, err := repo.GetFlexibleSavingsAccounts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, flexibleAccount.ID, accounts[0].ID)
	assert.Equal(t, flexibleAccount.AccountType, accounts[0].AccountType)
}

func TestAccountRepository_GetFlexibleSavingsAccounts_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	accounts, err := repo.GetFlexibleSavingsAccounts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, accounts)
}

func TestAccountRepository_UpdateBalance(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	repo := NewAccountRepository(db)

	// Setup test data
	user := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)

	testAccount := &account.Account{
		ID:            pkg.NewUUIDV7(),
		UserID:        user.ID,
		AccountType:   account.PaymentAccountType,
		AccountNumber: "PAY123456789",
		AccountName:   "Payment Account",
		Balance:       1000.0,
	}
	_, err = repo.Create(context.Background(), testAccount)
	require.NoError(t, err)

	newBalance := 2500.0
	err = repo.UpdateBalance(context.Background(), testAccount.ID, newBalance)

	assert.NoError(t, err)

	// Verify the balance was updated
	updatedAccount, err := repo.GetByID(context.Background(), testAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, newBalance, updatedAccount.Balance)
}

func TestAccountRepository_UpdateBalance_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	err := repo.UpdateBalance(context.Background(), pkg.NewUUIDV7(), 1000.0)

	assert.Error(t, err)
}