package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type BaseRepository[T any] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindById(ctx context.Context, id uuid.UUID) (T, error)
	Create(ctx context.Context, entity T) (T, error)
}
