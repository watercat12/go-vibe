package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/account"

	"github.com/labstack/echo/v4"
)

// CreatePaymentAccount godoc
//
//	@Summary		Create payment account
//	@Description	Create a new payment account for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	dto.AccountResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/accounts/payment [post]
//	@Security		BearerAuth
func (s *Server) CreatePaymentAccount(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	acc, err := s.AccountService.CreatePaymentAccount(c.Request().Context(), userID)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewAccountResponse(acc)
	return c.JSON(201, dto.Response{
		Status:  201,
		Message: "Payment account created successfully",
		Data:    resp,
	})
}

// CreateFixedSavingsAccount godoc
//
//	@Summary		Create fixed savings account
//	@Description	Create a new fixed-term savings account for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateFixedSavingsAccountRequest	true	"Fixed savings account creation data"
//	@Success		201		{object}	dto.AccountResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/accounts/savings/fixed [post]
//	@Security		BearerAuth
func (s *Server) CreateFixedSavingsAccount(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	var req dto.CreateFixedSavingsAccountRequest
	if err := c.Bind(&req); err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.BadRequestResponse)
	}

	domainReq := &account.CreateFixedSavingsAccountRequest{
		TermCode: req.TermCode,
	}

	acc, err := s.AccountService.CreateFixedSavingsAccount(c.Request().Context(), userID, domainReq)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewAccountResponse(acc)
	return c.JSON(201, dto.Response{
		Status:  201,
		Message: "Fixed savings account created successfully",
		Data:    resp,
	})
}

// CreateFlexibleSavingsAccount godoc
//
//	@Summary		Create flexible savings account
//	@Description	Create a new flexible savings account for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	dto.AccountResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/accounts/savings/flexible [post]
//	@Security		BearerAuth
func (s *Server) CreateFlexibleSavingsAccount(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	acc, err := s.AccountService.CreateFlexibleSavingsAccount(c.Request().Context(), userID)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewAccountResponse(acc)
	return c.JSON(201, dto.Response{
		Status:  201,
		Message: "Flexible savings account created successfully",
		Data:    resp,
	})
}

// ListAccounts godoc
//
//	@Summary		List user accounts
//	@Description	Get all accounts (payment and savings) for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.ListAccountsResponse
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/accounts [get]
//	@Security		BearerAuth
func (s *Server) ListAccounts(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	listResp, err := s.AccountService.ListAccounts(c.Request().Context(), userID)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	// Convert domain response to DTO
	var accounts []dto.AccountWithDetailsResponse
	for _, acc := range listResp.Accounts {
		dtoAcc := dto.AccountWithDetailsResponse{
			AccountResponse: dto.AccountResponse{
				ID:            acc.ID,
				UserID:        acc.UserID,
				AccountNumber: acc.AccountNumber,
				AccountType:   acc.AccountType,
				Balance:       acc.Balance,
				Status:        acc.Status,
				CreatedAt:     acc.CreatedAt,
				UpdatedAt:     acc.UpdatedAt,
			},
		}
		if acc.SavingsDetail != nil {
			dtoAcc.SavingsDetail = &dto.SavingsAccountDetailResponse{
				AccountID:             acc.SavingsDetail.AccountID,
				IsFixedTerm:           acc.SavingsDetail.IsFixedTerm,
				TermMonths:            acc.SavingsDetail.TermMonths,
				AnnualInterestRate:    acc.SavingsDetail.AnnualInterestRate,
				StartDate:             acc.SavingsDetail.StartDate,
				MaturityDate:          acc.SavingsDetail.MaturityDate,
				LastInterestCalcDate:  acc.SavingsDetail.LastInterestCalcDate,
			}
		}
		accounts = append(accounts, dtoAcc)
	}

	resp := dto.ListAccountsResponse{Accounts: accounts}
	return s.handleSuccess(c, resp)
}