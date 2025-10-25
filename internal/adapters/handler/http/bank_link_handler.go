package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) LinkBankAccount(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	var req dto.LinkBankAccountRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	linkedAccount, err := s.BankLinkService.LinkBankAccount(c.Request().Context(), userID, req.BankCode, req.AccountType)
	if err != nil {
		logrus.Error(err)
		return s.handleError(c, dto.BadRequestResponse)
	}

	resp := dto.NewLinkBankAccountResponse(linkedAccount)
	return c.JSON(http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    resp,
	})
}