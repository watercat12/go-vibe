package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"e-wallet/internal/adapters/handler/http/dto"
	"e-wallet/internal/config"
	"e-wallet/mocks"
	"e-wallet/pkg/logger"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		options     []Options
		expectError bool
	}{
		{
			name:        "success - no options",
			options:     nil,
			expectError: false,
		},
		{
			name: "success - with config option",
			options: []Options{
				WithConfig(&config.Config{AllowOrigins: "http://localhost:3000"}),
			},
			expectError: false,
		},
		{
			name: "error - option fails",
			options: []Options{
				func(s *Server) error {
					return assert.AnError
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := New(tt.options...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
				assert.NotNil(t, server.Router)
				assert.NotNil(t, server.Router.Validator)
			}
		})
	}
}

func TestServer_RegisterGlobalMiddlewares(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: &config.Config{AllowOrigins: "http://localhost:3000"},
		Logger: logger.NOOPLogger,
	}

	server.RegisterGlobalMiddlewares()

	// Check if middlewares are registered by checking routes or making a request
	// Since middlewares are internal, test by making a request to see if they work
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Test CORS if configured
	if server.Config.AllowOrigins != "" {
		req.Header.Set("Origin", "http://localhost:3000")
		server.Router.ServeHTTP(rec, req)
		// CORS headers should be set
		assert.Equal(t, "http://localhost:3000", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestServer_RegisterAuthMiddlewares(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	server.RegisterAuthMiddlewares()

	// Test by making a request to a protected route
	req := httptest.NewRequest(http.MethodPost, "/api/accounts/payment", nil)
	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, req)

	// Should return 400 bad request since no auth header (middleware returns error)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServer_RegisterUserClaimsMiddlewares(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	server.RegisterUserClaimsMiddlewares()

	// Test by making a request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, req)

	// Middleware should run without error
	assert.Equal(t, http.StatusNotFound, rec.Code) // No route defined yet
}

func TestServer_ServeHTTP(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	// Add a test route
	server.Router.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test", rec.Body.String())
}

func TestServer_RegisterHealthCheck(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	group := server.Router.Group("")
	server.RegisterHealthCheck(group)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	server.Router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "OK", response["status"])
	assert.Equal(t, "Service is up and running", response["message"])
}

func TestServer_handleError(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := server.Router.NewContext(req, rec)

	resErr := dto.Response{
		Status:  http.StatusBadRequest,
		Message: "Bad Request",
	}

	err := server.handleError(c, resErr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response dto.Response
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusBadRequest, response.Status)
	assert.Equal(t, "Bad Request", response.Message)
	assert.Nil(t, response.Data)
}

func TestServer_handleSuccess(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := server.Router.NewContext(req, rec)

	data := map[string]string{"key": "value"}

	err := server.handleSuccess(c, data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.Response
	_ = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, response.Status)
	assert.Equal(t, "OK", response.Message)
	assert.NotNil(t, response.Data)
}

func TestServer_requestID(t *testing.T) {
	server := &Server{
		Router: echo.New(),
		Config: config.Empty,
		Logger: logger.NOOPLogger,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderXRequestID, "test-request-id")
	rec := httptest.NewRecorder()
	c := server.Router.NewContext(req, rec)

	// Simulate middleware setting the header
	c.Response().Header().Set(echo.HeaderXRequestID, "test-request-id")

	requestID := server.requestID(c)

	assert.Equal(t, "test-request-id", requestID)
}

func TestServer_RegisterRoute(t *testing.T) {
	server := &Server{
		Router:         echo.New(),
		Config:         config.Empty,
		Logger:         logger.NOOPLogger,
		UserService:    mocks.NewMockUserService(t),
		AccountService: mocks.NewMockAccountService(t),
	}

	server.RegisterRoute()

	routes := server.Router.Routes()

	routePaths := make(map[string]bool)
	for _, route := range routes {
		routePaths[route.Path] = true
	}

	// Check if expected routes are registered
	expectedRoutes := []string{
		"/api/auth/register",
		"/api/auth/login",
		"/api/user/profile",
		"/api/accounts/payment",
		"/api/accounts/savings/fixed",
		"/api/accounts/savings/flexible",
	}

	for _, route := range expectedRoutes {
		assert.True(t, routePaths[route], "Route %s should be registered", route)
	}
}