package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/profile"

	"github.com/labstack/echo/v4"
)

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Update the authenticated user's profile information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UpdateProfileRequest	true	"Profile update data"
//	@Success		200		{object}	dto.ProfileResponse
//	@Failure		400		{object}	dto.Response
//	@Failure		401		{object}	dto.Response
//	@Failure		500		{object}	dto.Response
//	@Router			/api/users/profile [put]
//	@Security		BearerAuth
func (s *Server) UpdateProfile(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.BadRequestResponse)
	}

	if err := c.Validate(&req); err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.BadRequestResponse)
	}

	profileReq := &profile.UpdateProfileRequest{
		DisplayName: req.DisplayName,
		AvatarURL:   req.AvatarURL,
		PhoneNumber: req.PhoneNumber,
		NationalID:  req.NationalID,
		BirthYear:   req.BirthYear,
		Gender:      req.Gender,
		Team:        req.Team,
	}

	updatedProfile, err := s.ProfileService.UpdateProfile(c.Request().Context(), userID, profileReq)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewProfileResponse(updatedProfile)
	return s.handleSuccess(c, resp)
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieve the authenticated user's profile information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.ProfileResponse
//	@Failure		401	{object}	dto.Response
//	@Failure		500	{object}	dto.Response
//	@Router			/api/users/profile [get]
//	@Security		BearerAuth
func (s *Server) GetProfile(c echo.Context) error {
	userID := c.Get(UserIDKey).(string)
	if userID == "" {
		return s.handleError(c, dto.UnauthorizedResponse)
	}

	profile, err := s.ProfileService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		s.Logger.Error(err)
		return s.handleError(c, dto.InternalErrorResponse)
	}

	resp := dto.NewProfileResponse(profile)
	return s.handleSuccess(c, resp)
}