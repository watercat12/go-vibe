package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"e-wallet/internal/domain/bank_link"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg"
	_ "github.com/lib/pq"
)

func TestBankLinkRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)

	// Create a test user first
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	userRepo := NewUserRepository(db)
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	tests := []struct {
		name        string
		bankLink    *bank_link.BankLink
		expectError bool
	}{
		{
			name: "success - create bank link",
			bankLink: &bank_link.BankLink{
				ID:           pkg.NewUUIDV7(),
				UserID:       testUser.ID,
				BankCode:     "BANK001",
				AccountType:  "SAVINGS",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
				Status:       "ACTIVE",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.Create(context.Background(), tt.bankLink)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.bankLink.ID, result.ID)
				assert.Equal(t, tt.bankLink.UserID, result.UserID)
				assert.Equal(t, tt.bankLink.BankCode, result.BankCode)
				assert.Equal(t, tt.bankLink.AccountType, result.AccountType)
				assert.Equal(t, tt.bankLink.AccessToken, result.AccessToken)
				assert.Equal(t, tt.bankLink.RefreshToken, result.RefreshToken)
				assert.Equal(t, tt.bankLink.ExpiresIn, result.ExpiresIn)
				assert.Equal(t, tt.bankLink.Status, result.Status)
				assert.NotZero(t, result.CreatedAt)
			}
		})
	}
}

func TestBankLinkRepository_GetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)
	userRepo := NewUserRepository(db)

	// Create a test user first
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Setup test data
	testBankLink := &bank_link.BankLink{
		ID:           pkg.NewUUIDV7(),
		UserID:       testUser.ID,
		BankCode:     "BANK001",
		AccountType:  "SAVINGS",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
		Status:       "ACTIVE",
	}
	_, err = repo.Create(context.Background(), testBankLink)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectError bool
		expectedLen int
	}{
		{
			name:        "success - get bank links by user ID",
			userID:      testUser.ID,
			expectError: false,
			expectedLen: 1,
		},
		{
			name:        "success - no bank links found",
			userID:      pkg.NewUUIDV7(),
			expectError: false,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.GetByUserID(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
				if tt.expectedLen > 0 {
					assert.Equal(t, testBankLink.ID, result[0].ID)
					assert.Equal(t, testBankLink.UserID, result[0].UserID)
					assert.Equal(t, testBankLink.BankCode, result[0].BankCode)
					assert.Equal(t, testBankLink.AccountType, result[0].AccountType)
					assert.Equal(t, testBankLink.Status, result[0].Status)
				}
			}
		})
	}
}

func TestBankLinkRepository_CountByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)
	userRepo := NewUserRepository(db)

	// Create a test user first
	testUser := &user.User{
		ID:           pkg.NewUUIDV7(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	_, err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Setup test data
	testBankLink1 := &bank_link.BankLink{
		ID:           pkg.NewUUIDV7(),
		UserID:       testUser.ID,
		BankCode:     "BANK001",
		AccountType:  "SAVINGS",
		AccessToken:  "access-token-1",
		RefreshToken: "refresh-token-1",
		ExpiresIn:    3600,
		Status:       "ACTIVE",
	}
	testBankLink2 := &bank_link.BankLink{
		ID:           pkg.NewUUIDV7(),
		UserID:       testUser.ID,
		BankCode:     "BANK002",
		AccountType:  "CHECKING",
		AccessToken:  "access-token-2",
		RefreshToken: "refresh-token-2",
		ExpiresIn:    3600,
		Status:       "ACTIVE",
	}
	_, err = repo.Create(context.Background(), testBankLink1)
	require.NoError(t, err)
	_, err = repo.Create(context.Background(), testBankLink2)
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      string
		expectError bool
		expectedCount int
	}{
		{
			name:          "success - count bank links by user ID",
			userID:        testUser.ID,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "success - no bank links found",
			userID:        pkg.NewUUIDV7(),
			expectError:   false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.CountByUserID(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, result)
			}
		})
	}
}

func TestBankLinkRepository_Create_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	bankLink := &bank_link.BankLink{
		ID:           pkg.NewUUIDV7(),
		UserID:       "user-123",
		BankCode:     "BANK001",
		AccountType:  "SAVINGS",
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		ExpiresIn:    3600,
		Status:       "ACTIVE",
	}

	result, err := repo.Create(context.Background(), bankLink)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBankLinkRepository_GetByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.GetByUserID(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBankLinkRepository_CountByUserID_DBError(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBankLinkRepository(db)

	// Close the database connection to simulate DB error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	result, err := repo.CountByUserID(context.Background(), "user-123")

	assert.Error(t, err)
	assert.Equal(t, 0, result)
}