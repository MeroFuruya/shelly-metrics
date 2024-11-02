package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

func ValidateConnect(logger zerolog.Logger) func(context.Context, *pgconn.PgConn) error {
	return func(ctx context.Context, pgconn *pgconn.PgConn) error {
		address := pgconn.Conn().RemoteAddr().String()
		logger.Info().Str("address", address).Msg("Connected to database")
		return nil
	}
}

func Notice(logger zerolog.Logger) pgconn.NoticeHandler {
	return func(conn *pgconn.PgConn, notice *pgconn.Notice) {
		address := conn.Conn().RemoteAddr().String()
		logger.Warn().Str("address", address).
			Str("severity", notice.Severity).
			Str("code", notice.Code).
			Int32("position", notice.Position).
			Str("detail", notice.Detail).
			Str("schema", notice.SchemaName).
			Str("table", notice.TableName).
			Msg(notice.Message)
	}
}

func PgError(logger zerolog.Logger) pgconn.PgErrorHandler {
	return func(conn *pgconn.PgConn, err *pgconn.PgError) bool {
		address := conn.Conn().RemoteAddr().String()
		logger.Error().Str("address", address).
			Str("severity", err.Severity).
			Str("code", err.Code).
			Int32("position", err.Position).
			Str("detail", err.Detail).
			Str("schema", err.SchemaName).
			Str("table", err.TableName).
			Msg(err.Message)
		return true
	}
}
