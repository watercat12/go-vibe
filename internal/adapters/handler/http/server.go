package http

import (
	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/config"
	"e-wallet/internal/ports"
	"e-wallet/pkg/logger"
	"net/http"
	"strings"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "e-wallet/cmd/api/docs"
)

type Server struct {
	Router *echo.Echo
	Config *config.Config
	Logger *zap.SugaredLogger

	// service layers
	UserService    ports.UserService
	ProfileService ports.ProfileService
	AccountService ports.AccountService
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
	s.RegisterRoute()
	s.RegisterSwagger()

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
		"/swagger/",
	}
	s.Router.Use(NewAuthentication("header:Authorization", "Bearer", skipperPath).Middleware())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// RegisterHealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the service is up and running
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/healthz [get]
func (s *Server) RegisterHealthCheck(router *echo.Group) {
	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  http.StatusText(http.StatusOK),
			"message": "Service is up and running",
		})
	})
}

func (s *Server) handleError(c echo.Context, resErr dto.Response) error {
	s.Logger.Errorw(
		resErr.Message,
		zap.String("request_id", s.requestID(c)),
	)

	return c.JSON(resErr.Status, dto.Response{
		Status:  resErr.Status,
		Message: resErr.Message,
		Data:    nil,
	})
}

func (s *Server) handleSuccess(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
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

	// users
	apiGroup.PUT("/users/profile", s.UpdateProfile)
	apiGroup.GET("/users/profile", s.GetProfile)

	// accounts
	apiGroup.POST("/accounts/payment", s.CreatePaymentAccount)
	apiGroup.POST("/accounts/savings/fixed", s.CreateFixedSavingsAccount)
	apiGroup.POST("/accounts/savings/flexible", s.CreateFlexibleSavingsAccount)
	apiGroup.GET("/accounts", s.ListAccounts)
}

func (s *Server) RegisterSwagger() {
	s.Router.GET("/swagger/*", echoSwagger.WrapHandler)
}
