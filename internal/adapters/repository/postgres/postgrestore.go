package postgres

import (
	"e-wallet/internal/config"
	"fmt"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Options struct {
	DBName   string
	DBUser   string
	Password string
	Host     string
	Port     string
	SSLMode  bool
}

func ParseFromConfig(c *config.Config) Options {
	return Options{
		DBName:   c.DB.Name,
		DBUser:   c.DB.User,
		Password: c.DB.Pass,
		Host:     c.DB.Host,
		Port:     strconv.Itoa(c.DB.Port),
		SSLMode:  c.DB.EnableSSL,
	}
}

func NewConnection(opts Options) (*gorm.DB, error) {
	sslmode := "disable"
	if opts.SSLMode {
		sslmode = "enable"
	}

	datasource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		opts.Host, opts.Port, opts.DBUser, opts.Password, opts.DBName, sslmode,
	)

	db, err := gorm.Open(postgres.Open(datasource), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}