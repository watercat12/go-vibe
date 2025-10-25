package main

import (
	"fmt"
	"log"
	"net/http"

	"e-wallet/internal/adapters/external/bank_link"
	httpserver "e-wallet/internal/adapters/handler/http"
	"e-wallet/internal/adapters/repository/postgres"
	"e-wallet/internal/adapters/service"
	"e-wallet/internal/application/account"
	banklinkapp "e-wallet/internal/application/bank_link"
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
	repo := postgres.NewUserRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	ihRepo := postgres.NewInterestHistoryRepository(db)
	bankLinkRepo := postgres.NewBankLinkRepository(db)
	passwordService := service.NewBcryptPasswordService()
	bankLinkClient := bank_link.NewBankLinkClient("https://api.eazy-mock.teqn.asia")
	server.UserService = user.NewUserService(repo, profileRepo, passwordService)
	server.AccountService = account.NewAccountService(accountRepo, repo, profileRepo, txRepo, ihRepo)
	server.BankLinkService = banklinkapp.NewBankLinkService(bankLinkRepo, bankLinkClient)

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Info("server started!")
	applog.Fatal(http.ListenAndServe(addr, server))
}