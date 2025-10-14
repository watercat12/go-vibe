package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/user"

	"github.com/labstack/echo/v4"
)

func (s *Server) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	_, err := s.UserService.CreateUser(c.Request().Context(), &user.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewCreateUserResponse()
	return s.handleSuccess(c, resp)
}

func (s *Server) LoginUser(c echo.Context) error {
	var req dto.LoginUserRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	user, err := s.UserService.LoginUser(c.Request().Context(), &user.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	payload := TokenPayload{UserID: user.ID}
	token, err := CreateAccessToken(DefaultExpiredTime, payload, s.Config.JWTSecret)
	if err != nil {
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewLoginUserResponse(user, token)
	return s.handleSuccess(c, resp)
}

func (s *Server) UpdateProfile(c echo.Context) error {
	userID, ok := c.Get(UserIDKey).(string)
	if !ok {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		return s.handleError(c, dto.BadRequestResponse)
	}

	updatedProfile, err := s.UserService.UpdateProfile(c.Request().Context(), userID, &user.UpdateProfileRequest{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		PhoneNumber: req.PhoneNumber,
		NationalID:  req.NationalID,
		BirthYear:   req.BirthYear,
		Gender:      req.Gender,
		Team:        req.Team,
	})
	if err != nil {
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewUpdateProfileResponse(updatedProfile)
	return s.handleSuccess(c, resp)
}