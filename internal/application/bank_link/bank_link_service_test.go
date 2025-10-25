package bank_link

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"e-wallet/internal/domain/bank_link"
	"e-wallet/mocks"
)

func TestBankLinkService_LinkBankAccount(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		bankCode      string
		accountType   string
		mockSetup     func(*mocks.MockBankLinkRepository, *mocks.MockBankLinkClient)
		expectedError error
	}{
		{
			name:        "success - link bank account",
			userID:      "user-123",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			mockSetup: func(repo *mocks.MockBankLinkRepository, client *mocks.MockBankLinkClient) {
				repo.EXPECT().CountByUserID(mock.Anything, "user-123").Return(2, nil).Once()
				client.EXPECT().LinkAccount(mock.Anything, "BANK001", "SAVINGS").Return(&bank_link.BankLink{
					ID:           "link-123",
					UserID:       "",
					BankCode:     "BANK001",
					AccountType:  "SAVINGS",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					Status:       "ACTIVE",
				}, nil).Once()
				repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(bl *bank_link.BankLink) bool {
					return bl.UserID == "user-123" && bl.BankCode == "BANK001" && bl.AccountType == "SAVINGS"
				})).Return(&bank_link.BankLink{
					ID:           "link-123",
					UserID:       "user-123",
					BankCode:     "BANK001",
					AccountType:  "SAVINGS",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					Status:       "ACTIVE",
				}, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:        "error - count by user ID fails",
			userID:      "user-123",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			mockSetup: func(repo *mocks.MockBankLinkRepository, client *mocks.MockBankLinkClient) {
				repo.EXPECT().CountByUserID(mock.Anything, "user-123").Return(0, errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
		{
			name:        "error - max links exceeded",
			userID:      "user-123",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			mockSetup: func(repo *mocks.MockBankLinkRepository, client *mocks.MockBankLinkClient) {
				repo.EXPECT().CountByUserID(mock.Anything, "user-123").Return(5, nil).Once()
			},
			expectedError: errors.New("you have linked the maximum number of 5 bank accounts"),
		},
		{
			name:        "error - client link account fails",
			userID:      "user-123",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			mockSetup: func(repo *mocks.MockBankLinkRepository, client *mocks.MockBankLinkClient) {
				repo.EXPECT().CountByUserID(mock.Anything, "user-123").Return(2, nil).Once()
				client.EXPECT().LinkAccount(mock.Anything, "BANK001", "SAVINGS").Return(nil, errors.New("client error")).Once()
			},
			expectedError: errors.New("client error"),
		},
		{
			name:        "error - repository create fails",
			userID:      "user-123",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			mockSetup: func(repo *mocks.MockBankLinkRepository, client *mocks.MockBankLinkClient) {
				repo.EXPECT().CountByUserID(mock.Anything, "user-123").Return(2, nil).Once()
				client.EXPECT().LinkAccount(mock.Anything, "BANK001", "SAVINGS").Return(&bank_link.BankLink{
					ID:           "link-123",
					UserID:       "",
					BankCode:     "BANK001",
					AccountType:  "SAVINGS",
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    3600,
					Status:       "ACTIVE",
				}, nil).Once()
				repo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockBankLinkRepository(t)
			client := mocks.NewMockBankLinkClient(t)

			tt.mockSetup(repo, client)

			service := NewBankLinkService(repo, client)
			result, err := service.LinkBankAccount(context.Background(), tt.userID, tt.bankCode, tt.accountType)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.userID, result.UserID)
				assert.Equal(t, tt.bankCode, result.BankCode)
				assert.Equal(t, tt.accountType, result.AccountType)
				assert.Equal(t, "ACTIVE", result.Status)
			}
		})
	}
}