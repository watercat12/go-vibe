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
	"e-wallet/internal/domain/account"
	"e-wallet/mocks"
	"e-wallet/pkg/logger"
)

func TestServer_CreatePaymentAccount(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockSetup        func(*mocks.MockAccountService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name:   "success - create payment account",
			userID: "user-123",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreatePaymentAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&account.Account{
						ID:            "acc-123",
						UserID:        "user-123",
						AccountType:   account.PaymentAccountType,
						AccountNumber: "PAY123456789",
						Balance:       0.0,
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: dto.CreateAccountResponse{
					Account: &dto.AccountResponse{
						ID:            "acc-123",
						AccountNumber: "PAY123456789",
						Balance:       0.0,
					},
				},
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
		{
			name:   "error - service fails",
			userID: "user-123",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreatePaymentAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
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
			accountSvc := mocks.NewMockAccountService(t)

			// Setup mocks
			tt.mockSetup(accountSvc)

			// Create server
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/accounts/payment", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set user ID if provided
			if tt.userID != "" {
				c.Set(UserIDKey, tt.userID)
			}

			// Create server instance
			s := &Server{
				AccountService: accountSvc,
				Logger:         logger.NOOPLogger,
			}

			// Execute
			err := s.CreatePaymentAccount(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}

func TestServer_CreateFlexibleSavingsAccount(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		mockSetup        func(*mocks.MockAccountService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name:   "success - create flexible savings account",
			userID: "user-123",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreateFlexibleSavingsAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
					Return(&account.Account{
						ID:            "acc-123",
						UserID:        "user-123",
						AccountType:   account.FlexibleSavingsAccountType,
						AccountNumber: "SAV123456789",
						Balance:       0.0,
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: dto.CreateAccountResponse{
					Account: &dto.AccountResponse{
						ID:            "acc-123",
						AccountNumber: "SAV123456789",
						Balance:       0.0,
					},
				},
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
		{
			name:   "error - service fails",
			userID: "user-123",
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreateFlexibleSavingsAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					})).
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
			accountSvc := mocks.NewMockAccountService(t)

			// Setup mocks
			tt.mockSetup(accountSvc)

			// Create server
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/accounts/savings/flexible", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Set user ID if provided
			if tt.userID != "" {
				c.Set(UserIDKey, tt.userID)
			}

			// Create server instance
			s := &Server{
				AccountService: accountSvc,
				Logger:         logger.NOOPLogger,
			}

			// Execute
			err := s.CreateFlexibleSavingsAccount(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}

func TestServer_CreateFixedSavingsAccount(t *testing.T) {
	tests := []struct {
		name             string
		userID           string
		requestBody      dto.CreateFixedSavingsRequest
		mockSetup        func(*mocks.MockAccountService)
		expectedStatus   int
		expectedResponse dto.Response
	}{
		{
			name:   "success - create fixed savings account",
			userID: "user-123",
			requestBody: dto.CreateFixedSavingsRequest{
				TermMonths: 3,
			},
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreateFixedSavingsAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(termMonths int) bool {
						return termMonths == 3
					})).
					Return(&account.Account{
						ID:              "acc-123",
						UserID:          "user-123",
						AccountType:     account.FixedSavingsAccountType,
						AccountNumber:   "SAV123456789",
						Balance:         0.0,
						InterestRate:    func() *float64 { rate := 1.8; return &rate }(),
						FixedTermMonths: func() *int { term := 3; return &term }(),
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedResponse: dto.Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: dto.CreateAccountResponse{
					Account: &dto.AccountResponse{
						ID:              "acc-123",
						AccountNumber:   "SAV123456789",
						Balance:         0.0,
						InterestRate:    func() *float64 { rate := 1.8; return &rate }(),
						FixedTermMonths: func() *int { term := 3; return &term }(),
					},
				},
			},
		},
		{
			name:   "error - unauthorized",
			userID: "",
			requestBody: dto.CreateFixedSavingsRequest{
				TermMonths: 3,
			},
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedResponse: dto.UnauthorizedResponse,
		},
		{
			name:   "error - invalid request body",
			userID: "user-123",
			requestBody: dto.CreateFixedSavingsRequest{
				TermMonths: 3,
			},
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				// No mock setup needed as bind fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - validation fails",
			userID: "user-123",
			requestBody: dto.CreateFixedSavingsRequest{
				TermMonths: 2, // Invalid term
			},
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				// No mock setup needed as validation fails
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: dto.BadRequestResponse,
		},
		{
			name:   "error - service fails",
			userID: "user-123",
			requestBody: dto.CreateFixedSavingsRequest{
				TermMonths: 3,
			},
			mockSetup: func(accountSvc *mocks.MockAccountService) {
				accountSvc.EXPECT().
					CreateFixedSavingsAccount(mock.Anything, mock.MatchedBy(func(userID string) bool {
						return userID == "user-123"
					}), mock.MatchedBy(func(termMonths int) bool {
						return termMonths == 3
					})).
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
			accountSvc := mocks.NewMockAccountService(t)

			// Setup mocks
			tt.mockSetup(accountSvc)

			// Create server
			e := echo.New()
			var reqBody bytes.Buffer
			if tt.name != "error - invalid request body" {
				json.NewEncoder(&reqBody).Encode(tt.requestBody)
			} else {
				reqBody.WriteString("invalid json")
			}
			req := httptest.NewRequest(http.MethodPost, "/api/accounts/savings/fixed", &reqBody)
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
				AccountService: accountSvc,
				Logger:         logger.NOOPLogger,
			}

			// Execute
			err := s.CreateFixedSavingsAccount(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, rec.Code)
			expectedJSON, _ := json.Marshal(tt.expectedResponse)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}