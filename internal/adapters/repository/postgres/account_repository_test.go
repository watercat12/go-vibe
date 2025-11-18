package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"

	_ "github.com/lib/pq"
)


func TestAccountRepository_CreatePaymentAccount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	result, err := repo.CreatePaymentAccount(context.Background(), testUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.UserID)
	assert.Equal(t, "PAYMENT", result.AccountType)
	assert.Equal(t, 0.0, result.Balance)
	assert.Equal(t, "ACTIVE", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	assert.NotEmpty(t, result.AccountNumber)
	assert.Len(t, result.AccountNumber, 10)
}

func TestAccountRepository_CreatePaymentAccount_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	userID := pkg.NewUUIDV7()

	result, err := repo.CreatePaymentAccount(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAccountRepository_CreateFixedSavingsAccount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser2",
		Email:        "test2@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	req := &account.CreateFixedSavingsAccountRequest{
		TermCode: "12",
	}

	result, err := repo.CreateFixedSavingsAccount(context.Background(), testUser.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.UserID)
	assert.Equal(t, "FIXED_SAVINGS", result.AccountType)
	assert.Equal(t, 0.0, result.Balance)
	assert.Equal(t, "ACTIVE", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	assert.NotEmpty(t, result.AccountNumber)
	assert.Len(t, result.AccountNumber, 10)
}

func TestAccountRepository_CreateFixedSavingsAccount_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	userID := pkg.NewUUIDV7()
	req := &account.CreateFixedSavingsAccountRequest{
		TermCode: "12",
	}

	result, err := repo.CreateFixedSavingsAccount(context.Background(), userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAccountRepository_CreateFlexibleSavingsAccount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser3",
		Email:        "test3@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	result, err := repo.CreateFlexibleSavingsAccount(context.Background(), testUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.UserID)
	assert.Equal(t, "FLEXIBLE_SAVINGS", result.AccountType)
	assert.Equal(t, 0.0, result.Balance)
	assert.Equal(t, "ACTIVE", result.Status)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)
	assert.NotEmpty(t, result.AccountNumber)
	assert.Len(t, result.AccountNumber, 10)
}

func TestAccountRepository_CreateFlexibleSavingsAccount_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	userID := pkg.NewUUIDV7()

	result, err := repo.CreateFlexibleSavingsAccount(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAccountRepository_GetAccountsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser4",
		Email:        "test4@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test accounts
	_, err = repo.CreatePaymentAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
	_, err = repo.CreateFlexibleSavingsAccount(context.Background(), testUser.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectCount int
	}{
		{
			name:        "success - get accounts for user with accounts",
			userID:      testUser.ID,
			expectCount: 2,
		},
		{
			name:        "success - get accounts for user with no accounts",
			userID:      pkg.NewUUIDV7(),
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.GetAccountsByUserID(context.Background(), tt.userID)

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectCount)
			if tt.expectCount > 0 {
				for _, result := range results {
					assert.Equal(t, tt.userID, result.UserID)
				}
			}
		})
	}
}

func TestAccountRepository_GetAccountsByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	results, err := repo.GetAccountsByUserID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestAccountRepository_GetAccountByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser5",
		Email:        "test5@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test account
	testAccount, err := repo.CreatePaymentAccount(context.Background(), testUser.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		accountID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing account",
			accountID:   testAccount.ID,
			expectError: false,
		},
		{
			name:        "error - account not found",
			accountID:   pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetAccountByID(context.Background(), tt.accountID)

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
				assert.Equal(t, testAccount.AccountNumber, result.AccountNumber)
				assert.Equal(t, testAccount.AccountType, result.AccountType)
				assert.Equal(t, testAccount.Balance, result.Balance)
				assert.Equal(t, testAccount.Status, result.Status)
			}
		})
	}
}

func TestAccountRepository_GetAccountByID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetAccountByID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrUserNotFound, err)
}

func TestAccountRepository_CountPaymentAccountsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create test users first
	userRepo := NewUserRepository(db)
	testUser1 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser6",
		Email:        "test6@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser1)
	require.NoError(t, err)

	testUser2 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser7",
		Email:        "test7@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err = userRepo.Create(context.Background(), testUser2)
	require.NoError(t, err)

	// Create test accounts
	_, err = repo.CreatePaymentAccount(context.Background(), testUser1.ID)
	require.NoError(t, err)
	_, err = repo.CreatePaymentAccount(context.Background(), testUser1.ID)
	require.NoError(t, err)
	_, err = repo.CreateFlexibleSavingsAccount(context.Background(), testUser1.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectCount int64
	}{
		{
			name:        "success - count payment accounts for user with accounts",
			userID:      testUser1.ID,
			expectCount: 2,
		},
		{
			name:        "success - count payment accounts for user with no payment accounts",
			userID:      testUser2.ID,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := repo.CountPaymentAccountsByUserID(context.Background(), tt.userID)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectCount, count)
		})
	}
}

func TestAccountRepository_CountPaymentAccountsByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	count, err := repo.CountPaymentAccountsByUserID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAccountRepository_CountSavingsAccountsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create test users first
	userRepo := NewUserRepository(db)
	testUser1 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser8",
		Email:        "test8@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser1)
	require.NoError(t, err)

	testUser2 := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser9",
		Email:        "test9@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err = userRepo.Create(context.Background(), testUser2)
	require.NoError(t, err)

	// Create test accounts
	_, err = repo.CreateFixedSavingsAccount(context.Background(), testUser1.ID, &account.CreateFixedSavingsAccountRequest{TermCode: "12"})
	require.NoError(t, err)
	_, err = repo.CreateFlexibleSavingsAccount(context.Background(), testUser1.ID)
	require.NoError(t, err)
	_, err = repo.CreatePaymentAccount(context.Background(), testUser1.ID)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectCount int64
	}{
		{
			name:        "success - count savings accounts for user with accounts",
			userID:      testUser1.ID,
			expectCount: 2,
		},
		{
			name:        "success - count savings accounts for user with no savings accounts",
			userID:      testUser2.ID,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := repo.CountSavingsAccountsByUserID(context.Background(), tt.userID)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectCount, count)
		})
	}
}

func TestAccountRepository_CountSavingsAccountsByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	count, err := repo.CountSavingsAccountsByUserID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Equal(t, int64(0), count)
}

func TestAccountRepository_UpdateAccountBalance(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser10",
		Email:        "test10@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test account
	testAccount, err := repo.CreatePaymentAccount(context.Background(), testUser.ID)
	require.NoError(t, err)

	newBalance := 100.50

	err = repo.UpdateAccountBalance(context.Background(), testAccount.ID, newBalance)

	assert.NoError(t, err)

	// Verify the balance was updated
	updatedAccount, err := repo.GetAccountByID(context.Background(), testAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, newBalance, updatedAccount.Balance)
}

func TestAccountRepository_UpdateAccountBalance_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	err := repo.UpdateAccountBalance(context.Background(), pkg.NewUUIDV7(), 100.50)

	assert.Error(t, err)
}

func TestSavingsAccountDetailRepository_CreateSavingsAccountDetail(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := NewAccountRepository(db)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser11",
		Email:        "test11@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test account
	testAccount, err := accountRepo.CreateFixedSavingsAccount(context.Background(), testUser.ID, &account.CreateFixedSavingsAccountRequest{TermCode: "12"})
	require.NoError(t, err)

	termMonths := 12
	maturityDate := time.Now().AddDate(0, 12, 0)
	detail := &account.SavingsAccountDetail{
		AccountID:             testAccount.ID,
		IsFixedTerm:           true,
		TermMonths:            &termMonths,
		AnnualInterestRate:    5.0,
		StartDate:             time.Now(),
		MaturityDate:          &maturityDate,
		LastInterestCalcDate:  nil,
	}

	err = detailRepo.CreateSavingsAccountDetail(context.Background(), detail)

	assert.NoError(t, err)
}

func TestSavingsAccountDetailRepository_CreateSavingsAccountDetail_DBError(t *testing.T) {
	db := setupTestDB(t)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	detail := &account.SavingsAccountDetail{
		AccountID:             pkg.NewUUIDV7(),
		IsFixedTerm:           true,
		TermMonths:            nil,
		AnnualInterestRate:    5.0,
		StartDate:             time.Now(),
		MaturityDate:          nil,
		LastInterestCalcDate:  nil,
	}

	err := detailRepo.CreateSavingsAccountDetail(context.Background(), detail)

	assert.Error(t, err)
}

func TestSavingsAccountDetailRepository_GetSavingsAccountDetailByAccountID(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := NewAccountRepository(db)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser12",
		Email:        "test12@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test account and detail
	testAccount, err := accountRepo.CreateFlexibleSavingsAccount(context.Background(), testUser.ID)
	require.NoError(t, err)

	detail := &account.SavingsAccountDetail{
		AccountID:             testAccount.ID,
		IsFixedTerm:           false,
		TermMonths:            nil,
		AnnualInterestRate:    4.0,
		StartDate:             time.Now(),
		MaturityDate:          nil,
		LastInterestCalcDate:  nil,
	}
	err = detailRepo.CreateSavingsAccountDetail(context.Background(), detail)
	require.NoError(t, err)

	tests := []struct {
		name        string
		accountID   string
		expectError bool
		expectedErr error
	}{
		{
			name:        "success - get existing detail",
			accountID:   testAccount.ID,
			expectError: false,
		},
		{
			name:        "error - detail not found",
			accountID:   pkg.NewUUIDV7(),
			expectError: true,
			expectedErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detailRepo.GetSavingsAccountDetailByAccountID(context.Background(), tt.accountID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErr != nil {
					assert.Equal(t, tt.expectedErr, err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, detail.AccountID, result.AccountID)
				assert.Equal(t, detail.IsFixedTerm, result.IsFixedTerm)
				assert.Equal(t, detail.AnnualInterestRate, result.AnnualInterestRate)
			}
		})
	}
}

func TestSavingsAccountDetailRepository_GetSavingsAccountDetailByAccountID_DBError(t *testing.T) {
	db := setupTestDB(t)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := detailRepo.GetSavingsAccountDetailByAccountID(context.Background(), pkg.NewUUIDV7())

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NotEqual(t, ErrUserNotFound, err)
}

func TestSavingsAccountDetailRepository_UpdateLastInterestCalcDate(t *testing.T) {
	db := setupTestDB(t)
	accountRepo := NewAccountRepository(db)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Create a test user first
	userRepo := NewUserRepository(db)
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser13",
		Email:        "test13@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create test account and detail
	testAccount, err := accountRepo.CreateFlexibleSavingsAccount(context.Background(), testUser.ID)
	require.NoError(t, err)

	detail := &account.SavingsAccountDetail{
		AccountID:             testAccount.ID,
		IsFixedTerm:           false,
		TermMonths:            nil,
		AnnualInterestRate:    4.0,
		StartDate:             time.Now(),
		MaturityDate:          nil,
		LastInterestCalcDate:  nil,
	}
	err = detailRepo.CreateSavingsAccountDetail(context.Background(), detail)
	require.NoError(t, err)

	newDate := time.Now()

	err = detailRepo.UpdateLastInterestCalcDate(context.Background(), testAccount.ID, &newDate)

	assert.NoError(t, err)

	// Verify the date was updated
	updatedDetail, err := detailRepo.GetSavingsAccountDetailByAccountID(context.Background(), testAccount.ID)
	require.NoError(t, err)
	assert.NotNil(t, updatedDetail.LastInterestCalcDate)
	assert.Equal(t, newDate.Format("2006-01-02"), updatedDetail.LastInterestCalcDate.Format("2006-01-02"))
}

func TestSavingsAccountDetailRepository_UpdateLastInterestCalcDate_DBError(t *testing.T) {
	db := setupTestDB(t)
	detailRepo := NewSavingsAccountDetailRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	newDate := time.Now()

	err := detailRepo.UpdateLastInterestCalcDate(context.Background(), pkg.NewUUIDV7(), &newDate)

	assert.Error(t, err)
}