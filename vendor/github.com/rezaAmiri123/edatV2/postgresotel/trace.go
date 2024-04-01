package postgresotel

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgconn"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/stackus/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tracedDB struct {
	postgres.DB
}

var _ postgres.DB = (*tracedDB)(nil)

func Trace(db postgres.DB) postgres.DB {
	return tracedDB{DB: db}
}

func (t tracedDB) PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("PrepareContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", float64(time.Since(started).Seconds())),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.DB.PrepareContext(ctx, query)
}

func (t tracedDB) ExecContext(ctx context.Context, query string, args ...any) (result sql.Result, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("ExecContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", float64(time.Since(started).Seconds())),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.DB.ExecContext(ctx, query, args...)
}

func (t tracedDB) QueryContext(ctx context.Context, query string, args ...any) (rows *sql.Rows, err error) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", float64(time.Since(started).Seconds())),
		))
		t.recordError(span, err)
	}(time.Now())

	return t.DB.QueryContext(ctx, query, args...)
}

func (t tracedDB) QueryRowContext(ctx context.Context, query string, args ...any) (row *sql.Row) {
	span := trace.SpanFromContext(ctx)

	defer func(started time.Time) {
		span.AddEvent("QueryRowContext", trace.WithAttributes(
			attribute.String("Query", query),
			attribute.Float64("Took", float64(time.Since(started).Seconds())),
		))
		t.recordError(span, row.Err())
	}(time.Now())

	return t.DB.QueryRowContext(ctx, query, args...)
}

func (t tracedDB) recordError(span trace.Span, err error) {
	if err == nil {
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		span.AddEvent("Database Error", trace.WithAttributes(
			attribute.String("Error", err.Error()),
			attribute.String("Code", pgErr.Code),
			attribute.String("Severity", pgErr.Severity),
			attribute.String("Message", pgErr.Message),
			attribute.String("Detail", pgErr.Detail),
		))
	} else {
		span.AddEvent("Database Error", trace.WithAttributes(
			attribute.String("Error", err.Error()),
		))
	}
}
