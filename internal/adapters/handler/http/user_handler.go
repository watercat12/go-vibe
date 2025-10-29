package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/user"

	"github.com/labstack/echo/v4"
)

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Register a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"User registration data"
//	@Success		200		{object}	dto.CreateUserResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/auth/register [post]
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

// LoginUser godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginUserRequest	true	"User login data"
//	@Success		200		{object}	dto.LoginUserResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/auth/login [post]
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
