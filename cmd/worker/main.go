package main

import (
	"context"
	"log"

	"e-wallet/internal/adapters/repository/postgres"
	"e-wallet/internal/application/account"
	"e-wallet/internal/config"
	"e-wallet/pkg/logger"

	sentrygo "github.com/getsentry/sentry-go"
)

func main() {
	applog, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("cannot load config: %v\n", err)
	}

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

	repo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	ihRepo := postgres.NewInterestHistoryRepository(db)
	accountService := account.NewAccountService(accountRepo, repo, profileRepo, txRepo, ihRepo)

	// Run interest calculation
	if err := accountService.CalculateDailyInterest(context.Background()); err != nil {
		applog.Fatal(err)
	}

	applog.Info("Interest calculation completed successfully")
}