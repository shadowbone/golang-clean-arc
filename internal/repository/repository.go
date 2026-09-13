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
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type BaseRepository[T any] interface {
	FindAll(ctx context.Context, p Pagination) (Page[T], error)
	FindById(ctx context.Context, id uuid.UUID) (T, error)
	Create(ctx context.Context, entity T) (T, error)
	Update(ctx context.Context, id uuid.UUID, entity T) (T, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
