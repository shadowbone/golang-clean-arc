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
	Update(ctx context.Context, userId, id uuid.UUID, entity Address) (Address, error)
	Delete(ctx context.Context, userId, id uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Address, error)
	UnsetDefault(ctx context.Context, userId uuid.UUID) error
	SetDefault(ctx context.Context, userId, id uuid.UUID) error
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

func (d *DBAddress) Update(ctx context.Context, userId, id uuid.UUID, entity Address) (Address, error) {
	var a Address
	err := d.Q(ctx).QueryRow(ctx, `
		UPDATE addresses
		SET label = $3, street = $4, city = $5, updated_at = now()
		WHERE user_id = $1 and id = $2
		RETURNING id, user_id, label, street, city,created_at, updated_at
	`, userId, id, entity.Label, entity.Street, entity.City).Scan(
		&a.ID, &a.UserID, &a.Label, &a.Street, &a.City, &a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Address{}, repository.ErrNotFound
		}

		return Address{}, repository.ErrDuplicateKey
	}

	return a, nil
}

func (d *DBAddress) Delete(ctx context.Context, userId, id uuid.UUID) error {
	tag, err := d.Q(ctx).Exec(ctx, `
		DELETE FROM addresses where user_id = $1 and id = $2
	`, userId, id)

	if err != nil {
		return fmt.Errorf("delete addresses : %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (d *DBAddress) UnsetDefault(ctx context.Context, userId uuid.UUID) error {
	result, err := d.Q(ctx).Exec(ctx, `
		UPDATE addresses SET is_default = false, updated_at = now()
		WHERE user_id = $1`, userId)

	if err != nil {
		return fmt.Errorf("reset default addresses : %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (d *DBAddress) SetDefault(ctx context.Context, userId, id uuid.UUID) error {
	result, err := d.Q(ctx).Exec(ctx, `
		UPDATE addresses SET is_default = true, updated_at = now()
		WHERE id = $1 AND user_id = $2
	`, id, userId)

	if err != nil {
		return fmt.Errorf("set address default : %w", err)
	}

	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
