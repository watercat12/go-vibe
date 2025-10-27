package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/domain/profile"

	"github.com/labstack/echo/v4"
)

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