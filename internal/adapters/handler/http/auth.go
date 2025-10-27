package http

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const (
	UserClaimKey = "UserClaimKey"
 	UserIDKey = "UserID"

)

type Authentication struct {
	SkipperPath []string
	KeyLookup   string
	AuthScheme  string
}

func NewAuthentication(keyLookup string, authScheme string, skipperPath []string) *Authentication {
	return &Authentication{
		SkipperPath: skipperPath,
		KeyLookup:   keyLookup,
		AuthScheme:  authScheme,
	}
}

func (a *Authentication) Middleware() echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		return ContainFirst(a.SkipperPath, c.Path())
	}
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:    skipper,
		KeyLookup:  a.KeyLookup,
		AuthScheme: a.AuthScheme,
		Validator:  a.ValidateAccessToken,
	})
}

func (a *Authentication) ValidateAccessToken(token string, c echo.Context) (bool, error) {

	if token == "" {
		return false, errors.New("")
	}

	claims, err := ValidateToken(token, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return false, err
	}

	// check expired time
	now := time.Now().Local().Unix()
	if int64(claims["exp"].(float64)) < now {
		logrus.Error("Token is expired - ValidateAccessToken")

		return false, errors.New("token expired")
	}

	// get user_id
	payload, err := DecodeToken(claims)
	if err != nil {
		logrus.Error(err)

		return false, err
	}

	if payload.UserID == "" {
		logrus.Error("Unauthorized")

		return false, errors.New("")
	}

	c.Set(UserClaimKey, payload)
	c.Set(UserIDKey, payload.UserID)

	return true, nil
}

func ContainFirst(elems []string, v string) bool {
	for _, s := range elems {
		if strings.HasPrefix(v, s) {
			return true
		}
	}

	return false
}


