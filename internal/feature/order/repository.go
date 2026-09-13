package order

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type OrderRepository interface {
	Create(ctx context.Context, o Order) (Order, error)
	FindByID(ctx context.Context, UserID, id uuid.UUID) (Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, p repository.Pagination) (repository.Page[Order], error)
	UpdateStatus(ctx context.Context, userID, id uuid.UUID, status string) error
}

type DBOrderRepository struct {
	repository.DB
}

func NewOrderRepository(pool repository.Querier) OrderRepository {
	return &DBOrderRepository{DB: repository.NewDB(pool)}
}

func (d *DBOrderRepository) Create(ctx context.Context, o Order) (Order, error) {
	const qOrder = `
		INSERT INTO orders (user_id, total, status)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, total, status, created_at, updated_at`

	var out Order
	err := d.Q(ctx).QueryRow(ctx, qOrder, o.UserId, o.Total, o.Status).Scan(
		&out.ID, &out.UserId, &out.Total, &out.Status, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return Order{}, repository.Validationf("user tidak ditemukan")
		}

		return Order{}, fmt.Errorf("insert order: %w", err)
	}

	batch := &pgx.Batch{}
	const qItem = `
		INSERT INTO order_items (order_id, product_id, nama, qty, price, subtotal)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	for _, it := range o.Items {
		batch.Queue(qItem, out.ID, it.ProductID, it.Nama, it.Qty, it.Price, it.Subtotal)
	}

	br := d.Q(ctx).SendBatch(ctx, batch)

	out.Items = make([]OrderItem, 0, len(o.Items))
	for _, it := range o.Items {
		it.OrderID = out.ID
		if err := br.QueryRow().Scan(&it.ID); err != nil {
			_ = br.Close()
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return Order{}, repository.Validationf("produk tidak ditemukan")
			}
			return Order{}, fmt.Errorf("insert order item: %w", err)
		}
		out.Items = append(out.Items, it)
	}

	if err := br.Close(); err != nil {
		return Order{}, fmt.Errorf("close batch: %w", err)
	}

	return out, nil
}

func (d *DBOrderRepository) FindByID(ctx context.Context, UserID, id uuid.UUID) (Order, error) {
	const q = `
		SELECT id, user_id, total, status, created_at, updated_at
		FROM orders
		WHERE id = $1 and user_id = $2`
	var o Order
	err := d.Q(ctx).QueryRow(ctx, q, id, UserID).Scan(
		&o.ID, &o.UserId, &o.Total, &o.Status, &o.CreatedAt, &o.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Order{}, repository.ErrNotFound
		}
		return Order{}, fmt.Errorf("find order: %w", err)
	}

	items, err := d.findItems(ctx, []uuid.UUID{o.ID})
	if err != nil {
		return Order{}, err
	}

	o.Items = items[o.ID]

	return o, nil
}

func (d *DBOrderRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
	p repository.Pagination,
) (repository.Page[Order], error) {
	const countQ = `
		SELECT COUNT(*) FROM orders WHERE user_id = $1`

	var total int64
	if err := d.Q(ctx).QueryRow(ctx, countQ, userID).Scan(&total); err != nil {
		return repository.Page[Order]{}, fmt.Errorf("count order: %w", err)
	}

	const q = `
		SELECT id, user_id, total, status, created_at, update_at
		FROM orders
		WHERE user_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3`

	rows, err := d.Q(ctx).Query(ctx, q, p.Limit, p.Offset())
	if err != nil {
		return repository.Page[Order]{}, fmt.Errorf("query order: %w", err)
	}
	defer rows.Close()

	orders := make([]Order, 0, p.Limit)
	ids := make([]uuid.UUID, 0, p.Limit)

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.UserId, &o.Total, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return repository.Page[Order]{}, fmt.Errorf("scan order: %w", err)
		}
		orders = append(orders, o)
		ids = append(ids, o.ID)
	}

	if err := rows.Err(); err != nil {
		return repository.Page[Order]{}, fmt.Errorf("iterate orders: %w", err)
	}

	if len(ids) > 0 {
		items, err := d.findItems(ctx, ids)
		if err != nil {
			return repository.Page[Order]{}, err
		}

		for i := range orders {
			orders[i].Items = items[orders[i].ID]
		}
	}

	return repository.NewPage(orders, total, p), nil
}

func (d *DBOrderRepository) findItems(ctx context.Context, orderIDs []uuid.UUID) (map[uuid.UUID][]OrderItem, error) {
	const q = `
		SELECT id, order_id, product_id, nama, qty, price, subtotal
		FROM order_items
		WHERE order_id = ANY($1)
		ORDER BY id`

	rows, err := d.Q(ctx).Query(ctx, q, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]OrderItem)
	for rows.Next() {
		var it OrderItem
		if err := rows.Scan(&it.ID, &it.OrderID, &it.ProductID, &it.Nama, &it.Qty, &it.Price, &it.Subtotal); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}
		result[it.OrderID] = append(result[it.OrderID], it)
	}

	return result, rows.Err()
}

func (d *DBOrderRepository) UpdateStatus(ctx context.Context, userID, id uuid.UUID, status string) error {
	const q = `
		UPDATE orders
		SET status = $3, updated_at = now()
		WHERE id = $1 and user_id = $2`
	tag, err := d.Q(ctx).Exec(ctx, q, id, userID, status)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
