package database

import (
	"context"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTransactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *PgxTransactor {
	return &PgxTransactor{pool: pool}
}

func (t *PgxTransactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if repository.TxFromContext(ctx) != nil {
		return fn(ctx)
	}

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	if err := fn(repository.ContextWithTx(ctx, tx)); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
