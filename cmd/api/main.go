//	@title			E-Wallet API
//	@version		1.0
//	@description	E-Wallet service API
//	@host			pi.local:5111
//	@BasePath		/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
package main

import (
	"fmt"
	"log"
	"net/http"

	httpserver "e-wallet/internal/adapters/handler/http"
	"e-wallet/internal/adapters/repository/postgres"
	"e-wallet/internal/adapters/service"
	accountapp "e-wallet/internal/application/account"
	profileapp "e-wallet/internal/application/profile"
	"e-wallet/internal/application/user"
	"e-wallet/internal/config"
	"e-wallet/pkg/logger"

	sentrygo "github.com/getsentry/sentry-go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	applog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}
	defer logger.Sync(applog)

	cfg, err := config.LoadConfig()
	if err != nil {
		applog.Fatal(err)
	}

	err = sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.AppEnv,
		AttachStacktrace: true,
	})
	if err != nil {
		applog.Fatalf("cannot init sentry: %v", err)
	}

	db, err := postgres.NewConnection(postgres.ParseFromConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server, err := httpserver.New(httpserver.WithConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server.Logger = applog
	userRepo := postgres.NewUserRepository(db)
	passwordService := service.NewPasswordService()
	server.UserService = user.NewUserService(userRepo, passwordService)

	profileRepo := postgres.NewProfileRepository(db)
	server.ProfileService = profileapp.NewProfileService(userRepo, profileRepo)

	accountRepo := postgres.NewAccountRepository(db)
	savingsRepo := postgres.NewSavingsAccountDetailRepository(db)
	server.AccountService = accountapp.NewAccountService(userRepo, accountRepo, savingsRepo)

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Info("server started!")
	applog.Fatal(http.ListenAndServe(addr, server))
}