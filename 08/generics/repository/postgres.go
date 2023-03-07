package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type database struct {
	pool *pgxpool.Pool
}

func newDatabase(ctx context.Context, dsn string) (database, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return database{}, err
	}
	return database{
		pool: pool,
	}, nil
}

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func withConnection(ctx context.Context, p *pgxpool.Pool, f func(querier) error) (err error) {
	c, err := p.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer c.Release()
	return f(c)
}

func withConnectionResult[T any](ctx context.Context, p *pgxpool.Pool, f func(querier) (T, error)) (result T, err error) {
	c, err := p.Acquire(ctx)
	if err != nil {
		return result, fmt.Errorf("acquiring connection: %w", err)
	}
	defer c.Release()
	return f(c)
}

func withTransaction(ctx context.Context, p *pgxpool.Pool, f func(querier) error) (err error) {
	tx, err := p.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if rErr := tx.Rollback(ctx); rErr != nil && !errors.Is(rErr, pgx.ErrTxClosed) {
			err = errors.Join(rErr, err)
		}
	}()

	if err = f(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

func withTransactionResult[T any](ctx context.Context, s *pgxpool.Pool, f func(querier) (T, error)) (result T, err error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return result, fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() {
		if rErr := tx.Rollback(ctx); rErr != nil && !errors.Is(rErr, pgx.ErrTxClosed) {
			err = errors.Join(rErr, err)
		}
	}()

	result, err = f(tx)
	if err != nil {
		return result, err
	}

	if err = tx.Commit(ctx); err != nil {
		return result, fmt.Errorf("committing transaction: %w", err)
	}

	return result, nil
}

const (
	singleEntity = 1
)

func read[T any](ctx context.Context, querier querier, query string, args ...any) (*T, error) {
	var result T
	if err := pgxscan.Get(ctx, querier, &result, query, args...); err != nil {
		return nil, err
	}
	return &result, nil
}

func readValue[T any](ctx context.Context, querier querier, query string, args ...any) (T, error) {
	var result T
	if err := pgxscan.Get(ctx, querier, &result, query, args...); err != nil {
		return result, err
	}
	return result, nil
}

func list[T any](ctx context.Context, querier querier, query string, args ...any) ([]T, error) {
	var result []T
	if err := pgxscan.Select(ctx, querier, &result, query, args...); err != nil {
		return nil, err
	}
	return result, nil
}

func execOne(ctx context.Context, querier querier, query string, args ...any) error {
	r, err := querier.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if affected := r.RowsAffected(); affected != singleEntity {
		return fmt.Errorf("expected 1 row affected but affected %d rows", affected)
	}
	return nil
}

func exec(ctx context.Context, querier querier, query string, args ...any) error {
	if _, err := querier.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}
