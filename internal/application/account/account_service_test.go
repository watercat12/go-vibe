package account

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"e-wallet/internal/domain/account"
	"e-wallet/internal/domain/interest_history"
	"e-wallet/internal/domain/transaction"
	"e-wallet/internal/domain/user"
	"e-wallet/mocks"
)

func TestAccountService_CreatePaymentAccount(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mocks.MockAccountRepository, *mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedResult *account.Account
		expectedError  error
	}{
		{
			name:   "success - create payment account",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock no existing payment account
				accountRepo.EXPECT().
					GetPaymentAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("not found")).
					Once()

				// Mock account creation
				accountRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(acc *account.Account) bool {
					return acc.UserID == "user-123" &&
						acc.AccountType == account.PaymentAccountType &&
						acc.Balance == 0.0 &&
						len(acc.AccountNumber) > 0
				})).Return(&account.Account{
					ID:            "acc-123",
					UserID:        "user-123",
					AccountType:   account.PaymentAccountType,
					AccountNumber: "PAY123456789",
					Balance:       0.0,
				}, nil).Once()
			},
			expectedResult: &account.Account{
				ID:            "acc-123",
				UserID:        "user-123",
				AccountType:   account.PaymentAccountType,
				AccountNumber: "PAY123456789",
				Balance:       0.0,
			},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user not found
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("user not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:   "error - profile not found",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile not found
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("profile not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("profile not found"),
		},
		{
			name:   "error - payment account already exists",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock existing payment account
				accountRepo.EXPECT().
					GetPaymentAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&account.Account{
						ID:          "existing-acc-123",
						UserID:      "user-123",
						AccountType: account.PaymentAccountType,
					}, nil).Once()
			},
			expectedResult: nil,
			expectedError:  ErrLimitPaymentAccount,
		},
		{
			name:   "error - account creation fails",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock no existing payment account
				accountRepo.EXPECT().
					GetPaymentAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("not found")).
					Once()

				// Mock account creation fails
				accountRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil, errors.New("creation failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("creation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			accountRepo := mocks.NewMockAccountRepository(t)
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			// Setup mocks
			tt.mockSetup(accountRepo, userRepo, profileRepo)

			// Create service
			service := NewAccountService(accountRepo, userRepo, profileRepo, nil, nil)

			// Execute
			result, err := service.CreatePaymentAccount(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.UserID, result.UserID)
				assert.Equal(t, tt.expectedResult.AccountType, result.AccountType)
				assert.Equal(t, tt.expectedResult.Balance, result.Balance)
			}
		})
	}
}

func TestAccountService_CreateFlexibleSavingsAccount(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mocks.MockAccountRepository, *mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedResult *account.Account
		expectedError  error
	}{
		{
			name:   "success - create flexible savings account",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count < 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(2), nil).
					Once()

				// Mock account creation
				accountRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(acc *account.Account) bool {
					return acc.UserID == "user-123" &&
						acc.AccountType == account.FlexibleSavingsAccountType &&
						acc.Balance == 0.0 &&
						len(acc.AccountNumber) > 0
				})).Return(&account.Account{
					ID:            "acc-123",
					UserID:        "user-123",
					AccountType:   account.FlexibleSavingsAccountType,
					AccountNumber: "SAV123456789",
					Balance:       0.0,
				}, nil).Once()
			},
			expectedResult: &account.Account{
				ID:            "acc-123",
				UserID:        "user-123",
				AccountType:   account.FlexibleSavingsAccountType,
				AccountNumber: "SAV123456789",
				Balance:       0.0,
			},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user not found
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("user not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:   "error - profile not found",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile not found
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("profile not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("profile not found"),
		},
		{
			name:   "error - savings account limit reached",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count >= 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(5), nil).
					Once()
			},
			expectedResult: nil,
			expectedError:  ErrLimitSavingsAccount,
		},
		{
			name:   "error - count savings accounts fails",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock count fails
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(0), errors.New("count failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("count failed"),
		},
		{
			name:   "error - account creation fails",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count < 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(2), nil).
					Once()

				// Mock account creation fails
				accountRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil, errors.New("creation failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("creation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			accountRepo := mocks.NewMockAccountRepository(t)
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			// Setup mocks
			tt.mockSetup(accountRepo, userRepo, profileRepo)

			// Create service
			service := NewAccountService(accountRepo, userRepo, profileRepo, nil, nil)

			// Execute
			result, err := service.CreateFlexibleSavingsAccount(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.UserID, result.UserID)
				assert.Equal(t, tt.expectedResult.AccountType, result.AccountType)
				assert.Equal(t, tt.expectedResult.Balance, result.Balance)
			}
		})
	}
}

func TestAccountService_CreateFixedSavingsAccount(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		termMonths     int
		mockSetup      func(*mocks.MockAccountRepository, *mocks.MockUserRepository, *mocks.MockProfileRepository)
		expectedResult *account.Account
		expectedError  error
	}{
		{
			name:       "success - create fixed savings account",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count < 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(2), nil).
					Once()

				// Mock account creation
				accountRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(acc *account.Account) bool {
					return acc.UserID == "user-123" &&
						acc.AccountType == account.FixedSavingsAccountType &&
						acc.Balance == 0.0 &&
						*acc.FixedTermMonths == 3 &&
						*acc.InterestRate == 1.8 &&
						len(acc.AccountNumber) > 0
				})).Return(&account.Account{
					ID:              "acc-123",
					UserID:          "user-123",
					AccountType:     account.FixedSavingsAccountType,
					AccountNumber:   "SAV123456789",
					Balance:         0.0,
					InterestRate:    func() *float64 { rate := 1.8; return &rate }(),
					FixedTermMonths: func() *int { term := 3; return &term }(),
				}, nil).Once()
			},
			expectedResult: &account.Account{
				ID:              "acc-123",
				UserID:          "user-123",
				AccountType:     account.FixedSavingsAccountType,
				AccountNumber:   "SAV123456789",
				Balance:         0.0,
				InterestRate:    func() *float64 { rate := 1.8; return &rate }(),
				FixedTermMonths: func() *int { term := 3; return &term }(),
			},
			expectedError: nil,
		},
		{
			name:       "error - user not found",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user not found
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("user not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:       "error - profile not found",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile not found
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("profile not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("profile not found"),
		},
		{
			name:       "error - invalid term months",
			userID:     "user-123",
			termMonths: 2, // Invalid term
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()
			},
			expectedResult: nil,
			expectedError:  ErrInvalidTermMonths,
		},
		{
			name:       "error - savings account limit reached",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count >= 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(5), nil).
					Once()
			},
			expectedResult: nil,
			expectedError:  ErrLimitSavingsAccount,
		},
		{
			name:       "error - count savings accounts fails",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock count fails
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(0), errors.New("count failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("count failed"),
		},
		{
			name:       "error - account creation fails",
			userID:     "user-123",
			termMonths: 3,
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository, profileRepo *mocks.MockProfileRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock profile exists
				profileRepo.EXPECT().
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.Profile{ID: "profile-123", UserID: "user-123"}, nil).
					Once()

				// Mock savings accounts count < 5
				accountRepo.EXPECT().
					CountSavingsAccounts(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(int64(2), nil).
					Once()

				// Mock account creation fails
				accountRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil, errors.New("creation failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("creation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			accountRepo := mocks.NewMockAccountRepository(t)
			userRepo := mocks.NewMockUserRepository(t)
			profileRepo := mocks.NewMockProfileRepository(t)

			// Setup mocks
			tt.mockSetup(accountRepo, userRepo, profileRepo)

			// Create service
			service := NewAccountService(accountRepo, userRepo, profileRepo, nil, nil)

			// Execute
			result, err := service.CreateFixedSavingsAccount(context.Background(), tt.userID, tt.termMonths)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.UserID, result.UserID)
				assert.Equal(t, tt.expectedResult.AccountType, result.AccountType)
				assert.Equal(t, tt.expectedResult.Balance, result.Balance)
				assert.Equal(t, *tt.expectedResult.InterestRate, *result.InterestRate)
				assert.Equal(t, *tt.expectedResult.FixedTermMonths, *result.FixedTermMonths)
			}
		})
	}
}

func TestAccountService_CalculateDailyInterest(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.MockAccountRepository, *mocks.MockTransactionRepository, *mocks.MockInterestHistoryRepository)
		expectedError error
	}{
		{
			name: "success - calculate daily interest for multiple accounts",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, txRepo *mocks.MockTransactionRepository, ihRepo *mocks.MockInterestHistoryRepository) {
				// Mock GetFlexibleSavingsAccounts to return multiple accounts
				accounts := []*account.Account{
					{ID: "acc-1", UserID: "user-1", AccountType: account.FlexibleSavingsAccountType, Balance: 1000000, CreatedAt: time.Now().AddDate(0, 0, -90)},    // Older than 30 days, balance < 10M (annualRate = 0.003)
					{ID: "acc-2", UserID: "user-2", AccountType: account.FlexibleSavingsAccountType, Balance: 20000000, CreatedAt: time.Now().AddDate(0, 0, -10)}, // Newer than 30 days (annualRate = 0.008)
					{ID: "acc-3", UserID: "user-3", AccountType: account.FlexibleSavingsAccountType, Balance: 60000000, CreatedAt: time.Now().AddDate(0, 0, -90)}, // Older than 30 days, balance > 50M (annualRate = 0.005)
					{ID: "acc-4", UserID: "user-4", AccountType: account.FlexibleSavingsAccountType, Balance: 0, CreatedAt: time.Now().AddDate(0, 0, -90)},       // Zero balance (no interest)
					{ID: "acc-5", UserID: "user-5", AccountType: account.FlexibleSavingsAccountType, Balance: 30000000, CreatedAt: time.Now().AddDate(0, 0, -90)}, // Older than 30 days, 10M <= balance < 50M (annualRate = 0.004)
				}
				accountRepo.EXPECT().GetFlexibleSavingsAccounts(mock.Anything).Return(accounts, nil).Once()

				// Mock UpdateBalance for acc-1
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-1", mock.MatchedBy(func(balance float64) bool {
					return balance > 1000000 // Expecting interest to be added
				})).Return(nil).Once()

				// Mock Create transaction for acc-1
				txRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tx *transaction.Transaction) bool {
					return tx.AccountID == "acc-1" && tx.TransactionType == transaction.TransactionTypeInterest
				})).Return(&transaction.Transaction{}, nil).Once()

				// Mock Create interest history for acc-1
				ihRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(ih *interest_history.InterestHistory) bool {
					return ih.AccountID == "acc-1"
				})).Return(&interest_history.InterestHistory{}, nil).Once()

				// Mock UpdateBalance for acc-2
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-2", mock.MatchedBy(func(balance float64) bool {
					return balance > 20000000 // Expecting interest to be added
				})).Return(nil).Once()

				// Mock Create transaction for acc-2
				txRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tx *transaction.Transaction) bool {
					return tx.AccountID == "acc-2" && tx.TransactionType == transaction.TransactionTypeInterest
				})).Return(&transaction.Transaction{}, nil).Once()

				// Mock Create interest history for acc-2
				ihRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(ih *interest_history.InterestHistory) bool {
					return ih.AccountID == "acc-2"
				})).Return(&interest_history.InterestHistory{}, nil).Once()

				// Mock UpdateBalance for acc-3
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-3", mock.MatchedBy(func(balance float64) bool {
					return balance > 60000000 // Expecting interest to be added
				})).Return(nil).Once()

				// Mock Create transaction for acc-3
				txRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tx *transaction.Transaction) bool {
					return tx.AccountID == "acc-3" && tx.TransactionType == transaction.TransactionTypeInterest
				})).Return(&transaction.Transaction{}, nil).Once()

				// Mock Create interest history for acc-3
				ihRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(ih *interest_history.InterestHistory) bool {
					return ih.AccountID == "acc-3"
				})).Return(&interest_history.InterestHistory{}, nil).Once()

				// acc-4 (zero balance) should not trigger UpdateBalance, Create, Create

				// Mock UpdateBalance for acc-5
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-5", mock.MatchedBy(func(balance float64) bool {
					return balance > 30000000 // Expecting interest to be added
				})).Return(nil).Once()

				// Mock Create transaction for acc-5
				txRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tx *transaction.Transaction) bool {
					return tx.AccountID == "acc-5" && tx.TransactionType == transaction.TransactionTypeInterest
				})).Return(&transaction.Transaction{}, nil).Once()

				// Mock Create interest history for acc-5
				ihRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(ih *interest_history.InterestHistory) bool {
					return ih.AccountID == "acc-5"
				})).Return(&interest_history.InterestHistory{}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "error - GetFlexibleSavingsAccounts fails",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, txRepo *mocks.MockTransactionRepository, ihRepo *mocks.MockInterestHistoryRepository) {
				accountRepo.EXPECT().GetFlexibleSavingsAccounts(mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
		{
			name: "error - UpdateBalance fails",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, txRepo *mocks.MockTransactionRepository, ihRepo *mocks.MockInterestHistoryRepository) {
				accounts := []*account.Account{
					{ID: "acc-1", UserID: "user-1", AccountType: account.FlexibleSavingsAccountType, Balance: 1000000, CreatedAt: time.Now().AddDate(0, 0, -60)},
				}
				accountRepo.EXPECT().GetFlexibleSavingsAccounts(mock.Anything).Return(accounts, nil).Once()
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-1", mock.AnythingOfType("float64")).Return(errors.New("update error")).Once()
			},
			expectedError: errors.New("update error"),
		},
		{
			name: "error - Create transaction fails",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, txRepo *mocks.MockTransactionRepository, ihRepo *mocks.MockInterestHistoryRepository) {
				accounts := []*account.Account{
					{ID: "acc-1", UserID: "user-1", AccountType: account.FlexibleSavingsAccountType, Balance: 1000000, CreatedAt: time.Now().AddDate(0, 0, -60)},
				}
				accountRepo.EXPECT().GetFlexibleSavingsAccounts(mock.Anything).Return(accounts, nil).Once()
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-1", mock.AnythingOfType("float64")).Return(nil).Once()
				txRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*transaction.Transaction")).Return(nil, errors.New("tx error")).Once()
			},
			expectedError: errors.New("tx error"),
		},
		{
			name: "error - Create interest history fails",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, txRepo *mocks.MockTransactionRepository, ihRepo *mocks.MockInterestHistoryRepository) {
				accounts := []*account.Account{
					{ID: "acc-1", UserID: "user-1", AccountType: account.FlexibleSavingsAccountType, Balance: 1000000, CreatedAt: time.Now().AddDate(0, 0, -60)},
				}
				accountRepo.EXPECT().GetFlexibleSavingsAccounts(mock.Anything).Return(accounts, nil).Once()
				accountRepo.EXPECT().UpdateBalance(mock.Anything, "acc-1", mock.AnythingOfType("float64")).Return(nil).Once()
				txRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*transaction.Transaction")).Return(&transaction.Transaction{}, nil).Once()
				ihRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interest_history.InterestHistory")).Return(nil, errors.New("ih error")).Once()
			},
			expectedError: errors.New("ih error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			accountRepo := mocks.NewMockAccountRepository(t)
			txRepo := mocks.NewMockTransactionRepository(t)
			ihRepo := mocks.NewMockInterestHistoryRepository(t)

			// Setup mocks
			tt.mockSetup(accountRepo, txRepo, ihRepo)

			// Create service
			service := NewAccountService(accountRepo, nil, nil, txRepo, ihRepo)

			// Execute
			err := service.CalculateDailyInterest(context.Background())

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountService_GetAccountsByUserID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockSetup      func(*mocks.MockAccountRepository, *mocks.MockUserRepository)
		expectedResult []*account.Account
		expectedError  error
	}{
		{
			name:   "success - get accounts by user ID",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock get accounts
				accounts := []*account.Account{
					{ID: "acc-1", UserID: "user-123", AccountType: account.PaymentAccountType},
					{ID: "acc-2", UserID: "user-123", AccountType: account.FlexibleSavingsAccountType},
				}
				accountRepo.EXPECT().
					GetAccountsByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(accounts, nil).
					Once()
			},
			expectedResult: []*account.Account{
				{ID: "acc-1", UserID: "user-123", AccountType: account.PaymentAccountType},
				{ID: "acc-2", UserID: "user-123", AccountType: account.FlexibleSavingsAccountType},
			},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository) {
				// Mock user not found
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("user not found")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
		{
			name:   "error - get accounts fails",
			userID: "user-123",
			mockSetup: func(accountRepo *mocks.MockAccountRepository, userRepo *mocks.MockUserRepository) {
				// Mock user exists
				userRepo.EXPECT().
					GetByID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&user.User{ID: "user-123"}, nil).
					Once()

				// Mock get accounts fails
				accountRepo.EXPECT().
					GetAccountsByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(nil, errors.New("get accounts failed")).
					Once()
			},
			expectedResult: nil,
			expectedError:  errors.New("get accounts failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			accountRepo := mocks.NewMockAccountRepository(t)
			userRepo := mocks.NewMockUserRepository(t)

			// Setup mocks
			tt.mockSetup(accountRepo, userRepo)

			// Create service
			service := NewAccountService(accountRepo, userRepo, nil, nil, nil)

			// Execute
			result, err := service.GetAccountsByUserID(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.expectedResult))
				for i, acc := range result {
					assert.Equal(t, tt.expectedResult[i].ID, acc.ID)
					assert.Equal(t, tt.expectedResult[i].UserID, acc.UserID)
					assert.Equal(t, tt.expectedResult[i].AccountType, acc.AccountType)
				}
			}
		})
	}
}
