package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(context.Context) error) error
}

type transaction interface {
	Executor
	Commit() error
	Rollback() error
}

type beginTransaction func(context.Context, *sql.TxOptions) (transaction, error)

func (d *contextDatabase) WithinTransaction(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	if _, ok := d.transactionFrom(ctx); ok {
		return fn(ctx)
	}

	tx, err := d.begin(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if panicValue := recover(); panicValue != nil {
			_ = tx.Rollback()
			panic(panicValue)
		}
	}()

	txCtx := context.WithValue(ctx, d.key, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(
				err,
				fmt.Errorf("rollback transaction: %w", rollbackErr),
			)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
