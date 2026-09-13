package order

import (
	"context"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
)

type productAdapter struct {
	repository.DB
}

func NewProductAdapter(pool repository.Querier) ProductFinder {
	return &productAdapter{DB: repository.NewDB(pool)}
}

func (a *productAdapter) FindByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]ProductInfo, error) {
	const q = `SELECT id, name, price FROM products WHERE id = ANY($1)`

	rows, err := a.Q(ctx).Query(ctx, q, ids)
	if err != nil {
		return nil, fmt.Errorf("find products by ids: %w", err)
	}

	defer rows.Close()

	result := make(map[uuid.UUID]ProductInfo, len(ids))
	for rows.Next() {
		var p ProductInfo
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		result[p.ID] = p
	}

	return result, rows.Err()
}
