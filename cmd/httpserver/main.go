package main

import (
	"fmt"
	"log"
	"net/http"

	"e-wallet/adapters/httpserver"
	"e-wallet/adapters/postgrestore"
	"e-wallet/pkg/config"
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

	db, err := postgrestore.NewConnection(postgrestore.ParseFromConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server, err := httpserver.New(httpserver.WithConfig(cfg))
	if err != nil {
		applog.Fatal(err)
	}

	server.Logger = applog
	server.UserRepository = postgrestore.NewUserRepository(db)

	addr := fmt.Sprintf(":%d", cfg.Port)
	applog.Info("server started!")
	applog.Fatal(http.ListenAndServe(addr, server))
}