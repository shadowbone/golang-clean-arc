package order

import (
	"context"

	"github.com/google/uuid"
)

type ProductInfo struct {
	ID    uuid.UUID
	Nama  string
	Harga int64
}

type ProductFinder interface {
	FindByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ProductInfo, error)
}
