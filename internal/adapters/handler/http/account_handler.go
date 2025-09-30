package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreatePaymentAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	createdAccount, err := s.AccountService.CreatePaymentAccount(c.Request().Context(), userID)
	if err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) CreateFlexibleSavingsAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	createdAccount, err := s.AccountService.CreateFlexibleSavingsAccount(c.Request().Context(), userID)
	if err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) CreateFixedSavingsAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req dto.CreateFixedSavingsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdAccount, err := s.AccountService.CreateFixedSavingsAccount(c.Request().Context(), userID, req.TermMonths)
	if err != nil {
		return s.handleError(c, err, http.StatusBadRequest)
	}

	resp := dto.NewCreateAccountResponse(createdAccount)
	return c.JSON(http.StatusCreated, resp)
}