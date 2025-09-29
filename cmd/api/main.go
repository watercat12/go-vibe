package main

import (
	"fmt"
	"log"
	"net/http"

	httpserver "e-wallet/internal/adapters/handler/http"
	"e-wallet/internal/adapters/repository/postgres"
	"e-wallet/internal/application/account"
	"e-wallet/internal/application/user"
	"e-wallet/internal/config"
	"e-wallet/pkg/logger"

	sentrygo "github.com/getsentry/sentry-go"
)

func main() {
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
	repo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	server.UserService = user.NewUserService(repo, profileRepo)
	server.AccountService = account.NewAccountService(accountRepo, repo, profileRepo)

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Info("server started!")
	applog.Fatal(http.ListenAndServe(addr, server))
}