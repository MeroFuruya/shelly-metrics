package database

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DatabaseConfig struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
	Tls      bool
}

func ParseConfig(config DatabaseConfig) (*pgx.ConnConfig, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	pgConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	if !config.Tls {
		pgConfig.TLSConfig = nil
	}

	return pgConfig, nil
}
