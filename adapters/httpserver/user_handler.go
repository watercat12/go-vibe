package httpserver

import (
	"e-wallet/adapters/httpserver/model"
	"e-wallet/domain/user"
	"e-wallet/presenter"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateUser(c echo.Context) error {
	var req model.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	hashedPassword, err := user.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	user := &user.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: hashedPassword,
	}

	createdUser, err := s.UserRepository.Create(c.Request().Context(), user)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	payload := TokenPayload{UserID: strconv.Itoa(createdUser.ID)}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := presenter.NewCreateUserResponse(token)
	return c.JSON(http.StatusCreated, resp)
}