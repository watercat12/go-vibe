package httpserver

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdUser, err := s.UserService.CreateUser(c.Request().Context(), &user.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	payload := TokenPayload{UserID: createdUser.ID}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := dto.NewCreateUserResponse(token)
	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) LoginUser(c echo.Context) error {
	var req dto.LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := s.UserService.LoginUser(c.Request().Context(), &user.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return s.handleError(c, err, http.StatusUnauthorized)
	}

	payload := TokenPayload{UserID: user.ID}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := dto.NewLoginUserResponse(user, token)
	return c.JSON(http.StatusOK, resp)
}