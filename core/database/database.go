package database

import (
	"context"

	"github.com/MeroFuruya/shelly-metrics/core/logging"
	"github.com/jackc/pgx/v5"
)

var globalConnection *pgx.Conn

func InitConnection(ctx context.Context, config DatabaseConfig) error {
	logger := logging.GetLogger("database")

	configStruct, err := ParseConfig(config)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse database config")
		return err
	}

	configStruct.AfterConnect = ValidateConnect(logger)
	configStruct.OnPgError = PgError(logger)
	configStruct.OnNotice = Notice(logger)
	configStruct.Tracer = NewDatabaseTracer(logger)
	globalConnection, err = pgx.ConnectConfig(ctx, configStruct)
	return err
}

func GetConnection() *pgx.Conn {
	if globalConnection == nil {
		panic("Database not initialized")
	}
	return globalConnection
}
