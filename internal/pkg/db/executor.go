package db

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Database interface {
	Executor
	Transactor
}

type transactionKey struct {
	_ byte
}

type contextDatabase struct {
	fallback Executor
	begin    beginTransaction
	key      *transactionKey
}

func New(database *sql.DB) Database {
	if database == nil {
		panic("db: nil database")
	}

	return newContextDatabase(
		database,
		func(ctx context.Context, opts *sql.TxOptions) (transaction, error) {
			return database.BeginTx(ctx, opts)
		},
	)
}

func newContextDatabase(fallback Executor, begin beginTransaction) *contextDatabase {
	return &contextDatabase{
		fallback: fallback,
		begin:    begin,
		key:      &transactionKey{},
	}
}

func (d *contextDatabase) ExecContext(
	ctx context.Context,
	query string,
	args ...any,
) (sql.Result, error) {
	return d.executor(ctx).ExecContext(ctx, query, args...)
}

func (d *contextDatabase) QueryContext(
	ctx context.Context,
	query string,
	args ...any,
) (*sql.Rows, error) {
	return d.executor(ctx).QueryContext(ctx, query, args...)
}

func (d *contextDatabase) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) *sql.Row {
	return d.executor(ctx).QueryRowContext(ctx, query, args...)
}

func (d *contextDatabase) executor(ctx context.Context) Executor {
	tx, ok := d.transactionFrom(ctx)
	if ok {
		return tx
	}

	return d.fallback
}

func (d *contextDatabase) transactionFrom(ctx context.Context) (transaction, bool) {
	tx, ok := ctx.Value(d.key).(transaction)
	return tx, ok
}
