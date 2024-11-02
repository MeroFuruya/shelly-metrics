package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type TracerContextKey struct{}
type TracerContext struct {
	Sql     string
	Args    string
	Address string
}

type DatabaseTracer struct {
	logger zerolog.Logger
}

func argsToString(args []interface{}) string {
	var strArgs []string
	for _, arg := range args {
		strArgs = append(strArgs, fmt.Sprintf("%v", arg))
	}
	return strings.Join(strArgs, ", ")
}

func cleanupSql(sql string) string {
	rows := strings.Split(sql, "\n")
	var cleanedRows []string
	for _, row := range rows {
		cleanedRows = append(cleanedRows, strings.TrimSpace(row))
	}
	return strings.Join(cleanedRows, " ")
}

func (tracer *DatabaseTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracerContext := TracerContext{
		Sql:     cleanupSql(data.SQL),
		Args:    argsToString(data.Args),
		Address: conn.PgConn().Conn().RemoteAddr().String(),
	}
	return context.WithValue(ctx, TracerContextKey{}, tracerContext)
}

func (tracer *DatabaseTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	tracerContext := ctx.Value(TracerContextKey{}).(TracerContext)
	if data.Err != nil {
		tracer.logger.Debug().
			Str("address", tracerContext.Address).
			Str("args", tracerContext.Args).
			Msg(tracerContext.Sql)
	} else {
		tracer.logger.Error().
			Str("address", tracerContext.Address).
			Str("args", tracerContext.Args).
			Err(data.Err).
			Msg(tracerContext.Sql)
	}
}

func (tracer *DatabaseTracer) TraceBatchStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchStartData) context.Context {
	var tracerContext []TracerContext = []TracerContext{}
	for _, query := range data.Batch.QueuedQueries {
		tracerContext = append(tracerContext, TracerContext{
			Sql:     cleanupSql(query.SQL),
			Args:    argsToString(query.Arguments),
			Address: conn.PgConn().Conn().RemoteAddr().String(),
		})
	}
	return context.WithValue(ctx, TracerContextKey{}, tracerContext)
}

func (tracer *DatabaseTracer) TraceBatchQuery(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchQueryData) {
}
func (tracer *DatabaseTracer) TraceBatchEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceBatchEndData) {
	tracerContext := ctx.Value(TracerContextKey{}).([]TracerContext)
	if data.Err != nil {
		tracer.logger.Err(data.Err).Msg("Error in batch query")
		for _, context := range tracerContext {
			tracer.logger.Err(data.Err).
				Str("address", context.Address).
				Str("args", context.Args).
				Msg(context.Sql)
		}
	} else {
		tracer.logger.Debug().
			Msg("Batch query executed")
		for _, context := range tracerContext {
			tracer.logger.Debug().
				Str("address", context.Address).
				Str("args", context.Args).
				Msg(context.Sql)
		}
	}
}

func NewDatabaseTracer(logger zerolog.Logger) *DatabaseTracer {
	return &DatabaseTracer{logger: logger}
}
