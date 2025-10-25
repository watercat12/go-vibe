package bank_link

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBankLinkClient_LinkAccount(t *testing.T) {
	tests := []struct {
		name          string
		bankCode      string
		accountType   string
		serverSetup   func() *httptest.Server
		expectedError error
	}{
		{
			name:        "success - link account",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"access_token": "test-access-token",
						"bank_code": "BANK001",
						"expires_in": 3600,
						"refresh_token": "test-refresh-token"
					}`))
				}))
			},
			expectedError: nil,
		},
		{
			name:        "error - http request fails",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			expectedError: &BankLinkError{Message: "failed to link account"},
		},
		{
			name:        "error - invalid json response",
			bankCode:    "BANK001",
			accountType: "SAVINGS",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`invalid json`))
				}))
			},
			expectedError: errors.New("invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverSetup()
			defer server.Close()

			client := NewBankLinkClient(server.URL + "/mock/bbecc26a-c96a-4aa6-8178-0ee9cde9f390/303/bank-link")
			result, err := client.LinkAccount(context.Background(), tt.bankCode, tt.accountType)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.name == "error - http request fails" {
					var bankLinkErr *BankLinkError
					assert.ErrorAs(t, err, &bankLinkErr)
					assert.Equal(t, "failed to link account", bankLinkErr.Message)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.bankCode, result.BankCode)
				assert.Equal(t, tt.accountType, result.AccountType)
				assert.Equal(t, "test-access-tokenn", result.AccessToken)
				assert.Equal(t, "test-refresh-tokenn", result.RefreshToken)
				assert.Equal(t, 3600, result.ExpiresIn)
				assert.Equal(t, "ACTIVE", result.Status)
			}
		})
	}
}

func TestBankLinkClient_LinkAccount_NetworkError(t *testing.T) {
	client := NewBankLinkClient("http://invalid-url")
	result, err := client.LinkAccount(context.Background(), "BANK001", "SAVINGS")

	assert.Error(t, err)
	assert.Nil(t, result)
}