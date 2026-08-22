package repository

import "context"

type DB struct {
	pool Querier
}

func NewDB(pool Querier) DB {
	return DB{pool: pool}
}

func (d DB) Q(ctx context.Context) Querier {
	if tx := TxFromContext(ctx); tx != nil {
		return tx
	}

	return d.pool
}
