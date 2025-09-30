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