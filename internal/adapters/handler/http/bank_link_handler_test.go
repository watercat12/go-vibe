package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/config"
	"e-wallet/internal/domain/bank_link"
	"e-wallet/mocks"
	"e-wallet/pkg/logger"
)

func TestServer_LinkBankAccount(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		requestBody      dto.LinkBankAccountRequest
		mockSetup        func(*mocks.MockBankLinkService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name:   "success - link bank account",
			userID: "user-123",
			requestBody: dto.LinkBankAccountRequest{
				BankCode:    "BANK001",
				AccountType: "SAVINGS",
			},
			mockSetup: func(bankLinkSvc *mocks.MockBankLinkService) {
				bankLinkSvc.EXPECT().
					LinkBankAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), "BANK001", "SAVINGS").
					Return(&bank_link.BankLink{
						ID:          "link-123",
						UserID:      "user-123",
						BankCode:    "BANK001",
						AccountType: "SAVINGS",
						Status:      "ACTIVE",
					}, nil).
					Once()
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: dto.Response{
				Status:  http.StatusCreated,
				Message: http.StatusText(http.StatusCreated),
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			requestBody: dto.LinkBankAccountRequest{
				BankCode:    "BANK001",
				AccountType: "SAVINGS",
			},
			mockSetup: func(bankLinkSvc *mocks.MockBankLinkService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
		{
			name:   "error - invalid request body",
			userID: "user-123",
			requestBody: dto.LinkBankAccountRequest{
				BankCode:    "BANK001",
				AccountType: "SAVINGS",
			},
			mockSetup: func(bankLinkSvc *mocks.MockBankLinkService) {
				// No mock setup needed as bind fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - validation fails",
			userID: "user-123",
			requestBody: dto.LinkBankAccountRequest{
				BankCode:    "",
				AccountType: "",
			},
			mockSetup: func(bankLinkSvc *mocks.MockBankLinkService) {
				// No mock setup needed as validation fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - service fails",
			userID: "user-123",
			requestBody: dto.LinkBankAccountRequest{
				BankCode:    "BANK001",
				AccountType: "SAVINGS",
			},
			mockSetup: func(bankLinkSvc *mocks.MockBankLinkService) {
				bankLinkSvc.EXPECT().
					LinkBankAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), "BANK001", "SAVINGS").
					Return(nil, errors.New("service error")).
					Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			bankLinkSvc := mocks.NewMockBankLinkService(t)

			// Setup mocks
			tt.mockSetup(bankLinkSvc)

			// Create server
			e := echo.New()
			var reqBody bytes.Buffer
			if tt.name != "error - invalid request body" {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			} else {
				reqBody.WriteString("invalid json")
			}
			req := httptest.NewRequest(http.MethodPost, "/api/user/bank-accounts", &reqBody)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set validator for tests that reach validation
			if tt.name != "error - unauthorized" && tt.name != "error - invalid request body" {
				v := validator.New()
				dto.RegisterCustomValidations(v)
				e.Validator = &CustomValidator{validator: v}
			}

			// Set user ID if provided
			if tt.userID != "" {
				c.Set(UserIDKey, tt.userID)
			}

			// Create server instance
			s := &Server{
				BankLinkService: bankLinkSvc,
				Logger:          logger.NOOPLogger,
				Config:          &config.Config{},
			}

			// Execute
			err := s.LinkBankAccount(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			var actualResponse dto.Response
			_ = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			dataActualJson, _ := json.Marshal(actualResponse.Data)
			var dataActual dto.LinkBankAccountResponse
			_ = json.Unmarshal(dataActualJson, &dataActual)
			expectedResponseJson, _ := json.Marshal(tt.expectedResponse.Data)
			var expectedData dto.LinkBankAccountResponse
			_ = json.Unmarshal(expectedResponseJson, &expectedData)

			assert.Equal(t, dto.Response{
				Status:  tt.expectedResponse.Status,
				Message: tt.expectedResponse.Message,
			}, dto.Response{
				Status:  actualResponse.Status,
				Message: actualResponse.Message,
			})
			if tt.name == "success - link bank account" {
				assert.NotEmpty(t, dataActual.ID)
				assert.Equal(t, "BANK001", dataActual.BankCode)
				assert.Equal(t, "SAVINGS", dataActual.AccountType)
				assert.Equal(t, "ACTIVE", dataActual.Status)
			}
		})
	}
}