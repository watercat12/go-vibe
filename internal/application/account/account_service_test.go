package account

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"e-wallet/internal/domain/account"
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
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
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
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
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
					GetByUserID(mock.Anything, mock.MatchedBy(func(userID string) bool {
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
