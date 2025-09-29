package httpserver

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/config"
	"e-wallet/internal/domain/user"
	"e-wallet/pkg/logger"
	"e-wallet/pkg/sentry"
	"net/http"
	"strings"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Router *echo.Echo
	Config *config.Config
	Logger *zap.SugaredLogger

	// service layers
	UserService user.UserService
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func New(options ...Options) (*Server, error) {
	s := Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	v := validator.New()
	dto.RegisterCustomValidations(v)
	s.Router.Validator = &CustomValidator{validator: v}

	for _, fn := range options {
		if err := fn(&s); err != nil {
			return nil, err
		}
	}

	s.RegisterGlobalMiddlewares()
	s.RegisterAuthMiddlewares()
	s.RegisterUserClaimsMiddlewares()
	s.RegisterRoute()

	s.RegisterHealthCheck(s.Router.Group(""))

	return &s, nil
}

func (s *Server) RegisterGlobalMiddlewares() {
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.Secure())
	s.Router.Use(middleware.RequestID())
	s.Router.Use(middleware.Gzip())
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	// CORS
	if s.Config.AllowOrigins != "" {
		aos := strings.Split(s.Config.AllowOrigins, ",")
		s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: aos,
		}))
	}
}

func (s *Server) RegisterAuthMiddlewares() {
	skipperPath := []string{
		"/healthz",
		"/api/auth",
	}
	s.Router.Use(NewAuthentication("header:Authorization", "Bearer", skipperPath).Middleware())
}

func (s *Server) RegisterUserClaimsMiddlewares() {
	s.Router.Use(CheckUserTypeMiddleware())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) RegisterHealthCheck(router *echo.Group) {
	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  http.StatusText(http.StatusOK),
			"message": "Service is up and running",
		})
	})
}

func (s *Server) handleError(c echo.Context, err error, status int) error {
	s.Logger.Errorw(
		err.Error(),
		zap.String("request_id", s.requestID(c)),
	)

	if status >= http.StatusInternalServerError {
		sentry.WithContext(c).Error(err)
	}

	return c.JSON(status, map[string]string{
		"message": http.StatusText(status),
	})
}

func (s *Server) requestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

func (s *Server) RegisterRoute() {
	apiGroup := s.Router.Group("/api")
	// auth
	apiGroup.POST("/auth/register", s.CreateUser)
	apiGroup.POST("/auth/login", s.LoginUser)
}
