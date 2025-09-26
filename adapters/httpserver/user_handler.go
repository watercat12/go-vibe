package httpserver

import (
	"e-wallet/adapters/httpserver/model"
	userdomain "e-wallet/domain/user"
	"e-wallet/pkg"
	"e-wallet/presenter"
	"net/http"

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

	hashedPassword, err := userdomain.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	usr := &userdomain.User{
		ID:           pkg.NewUUIDV7(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.UserRepository.Create(c.Request().Context(), usr)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	payload := TokenPayload{UserID: createdUser.ID}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := presenter.NewCreateUserResponse(token)
	return c.JSON(http.StatusCreated, resp)
}

func (s *Server) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	usr, err := s.UserRepository.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	if err := userdomain.CheckPassword(usr.PasswordHash, req.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	payload := TokenPayload{UserID: usr.ID}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := presenter.NewLoginResponse(token)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) GetMe(c echo.Context) error {
	userIDStr, ok := c.Get(UserIDKey).(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user, err := s.UserRepository.GetByID(c.Request().Context(), userIDStr)
	if err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	resp := presenter.NewUserResponse(user)
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(c echo.Context) error {
	userIDStr, ok := c.Get(UserIDKey).(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req model.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	profile := &userdomain.UserProfile{
		UserID:    userIDStr,
		Name:      req.Name,
		Email:     req.Email,
		Avatar:    req.Avatar,
		Phone:     req.Phone,
		IDNumber:  req.IDNumber,
		BirthYear: req.BirthYear,
		Gender:    req.Gender,
		Team:      req.Team,
	}

	if err := s.UserRepository.UpdateProfile(c.Request().Context(), profile); err != nil {
		return s.handleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
