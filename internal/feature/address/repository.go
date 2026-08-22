package address

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AddressRepository interface {
	Create(ctx context.Context, a Address) (Address, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Address, error)
}

type DBAddress struct {
	repository.DB
}

func NewAddressRepository(pool repository.Querier) AddressRepository {
	return &DBAddress{DB: repository.NewDB(pool)}
}

func (d *DBAddress) Create(ctx context.Context, a Address) (Address, error) {
	const q = `
		INSERT INTO addresses(user_id, label, street,city,is_default)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, user_id, label,street,city,is_default,created_at,updated_at`
	var out Address
	err := d.Q(ctx).QueryRow(ctx, q, a.UserID, a.Label, a.Street, a.City, a.IsDefault).Scan(
		&out.ID, &out.UserID, &out.Label, &out.Street, &out.City,
		&out.IsDefault, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return Address{}, repository.Validationf("user tidak ditemukan")
		}

		return Address{}, fmt.Errorf("create address: %w", err)
	}
	return out, nil
}

func (d *DBAddress) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	const count = `SELECT COUNT(*) FROM addresses WHERE user_id = $1`
	var total int64
	if err := d.Q(ctx).QueryRow(ctx, count, userID).Scan(&total); err != nil {
		return total, fmt.Errorf("count by user id:%w", err)
	}
	return total, nil
}

func (d *DBAddress) FindByUserID(ctx context.Context, userID uuid.UUID) ([]Address, error) {
	const q = `SELECT * FROM addresses WHERE user_id = $1`
	rows, err := d.Q(ctx).Query(ctx, q, userID)
	if err != nil {
		return []Address{}, fmt.Errorf("query address FindByUserId: %w", err)
	}
	defer rows.Close()

	address, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Address])
	if err != nil {
		return []Address{}, fmt.Errorf("collect address by user :%w", err)
	}
	return address, nil
}
