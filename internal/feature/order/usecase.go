package order

import (
	"context"
	"golang-rest-api/internal/logger"
	"golang-rest-api/internal/repository"
	"golang-rest-api/internal/validation"
	"log/slog"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo    OrderRepository
	product ProductFinder
	tx      repository.Transactor
}

func NewOrderUseCase(
	repo OrderRepository,
	product ProductFinder,
	tx repository.Transactor,
) *OrderUseCase {
	return &OrderUseCase{repo: repo, product: product, tx: tx}
}

func (u *OrderUseCase) Create(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (Order, error) {
	if err := validation.Struct(req); err != nil {
		return Order{}, err
	}

	qtyByProduct := make(map[uuid.UUID]int, len(req.Items))
	ids := make([]uuid.UUID, 0, len(req.Items))

	for _, it := range req.Items {
		if _, ada := qtyByProduct[it.ProductId]; !ada {
			ids = append(ids, it.ProductId)
		}
		qtyByProduct[it.ProductId] += it.Qty
	}

	var created Order

	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		prices, err := u.product.FindByIDs(ctx, ids)
		if err != nil {
			return err
		}

		items := make([]OrderItem, 0, len(ids))
		for _, id := range ids {
			p, ok := prices[id]
			if !ok {
				return repository.Validationf("product %s tidak ditemukan", id)
			}
			items = append(items, OrderItem{
				ProductID: p.ID,
				Nama:      p.Nama,
				Qty:       qtyByProduct[id],
				Price:     p.Harga,
			})
		}

		o := Order{
			UserId: userID,
			Status: StatusPending,
			Items:  items,
		}

		o.HitungTotal()

		created, err = u.repo.Create(ctx, o)
		return err
	})

	if err != nil {
		return Order{}, err
	}

	logger.FromContext(ctx).Info("order dibuat",
		slog.String("order_id", created.ID.String()),
		slog.String("user_id", userID.String()),
		slog.Int64("total", created.Total),
		slog.Int("jumlah_item", len(created.Items)),
	)

	return created, nil
}

func (u *OrderUseCase) GetByID(ctx context.Context, userID, id uuid.UUID) (Order, error) {
	return u.repo.FindByID(ctx, userID, id)
}

func (u *OrderUseCase) List(ctx context.Context, userID uuid.UUID, p repository.Pagination) (repository.Page[Order], error) {
	p.Normalize()
	return u.repo.FindByUserID(ctx, userID, p)
}

func (u *OrderUseCase) Cancel(ctx context.Context, userID, id uuid.UUID) error {
	return u.tx.WithinTx(ctx, func(ctx context.Context) error {
		o, err := u.repo.FindByID(ctx, userID, id)
		if err != nil {
			return err
		}

		if !o.BolehDibatalkan() {
			return repository.Validationf("order dengan status %s tidak bisa dibatalkan", o.Status)
		}

		return u.repo.UpdateStatus(ctx, userID, id, StatusCancelled)
	})
}
