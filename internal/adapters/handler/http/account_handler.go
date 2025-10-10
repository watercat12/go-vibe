package http

import (
	"e-wallet/internal/adapters/handler/http/dto"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreatePaymentAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	createdAccount, err := s.AccountService.CreatePaymentAccount(c.Request().Context(), userID)
	if err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return s.handleSuccess(c, resp)
}

func (s *Server) CreateFlexibleSavingsAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	createdAccount, err := s.AccountService.CreateFlexibleSavingsAccount(c.Request().Context(), userID)
	if err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return s.handleSuccess(c, resp)
}

func (s *Server) CreateFixedSavingsAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	var req dto.CreateFixedSavingsRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	createdAccount, err := s.AccountService.CreateFixedSavingsAccount(c.Request().Context(), userID, req.TermMonths)
	if err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return s.handleSuccess(c, resp)
}
