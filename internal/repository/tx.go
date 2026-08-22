package repository

import "context"

type txKey struct{}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

func ContextWithTx(ctx context.Context, q Querier) context.Context {
	return context.WithValue(ctx, txKey{}, q)
}

func TxFromContext(ctx context.Context) Querier {
	q, _ := ctx.Value(txKey{}).(Querier)
	return q
}
