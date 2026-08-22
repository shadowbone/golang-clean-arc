package product

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	repository.BaseRepository[Product]
}

type DBProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &DBProductRepository{db: db}
}

func (d *DBProductRepository) FindAll(ctx context.Context) ([]Product, error) {
	const q = `
		select id, name, price, created_at, updated_at
		FROM products
		ORDER BY id DESC
		LIMIT 100`
	rows, err := d.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[Product])
	if err != nil {
		return nil, fmt.Errorf("collect products: %w", err)
	}

	return items, nil
}

func (d *DBProductRepository) FindById(ctx context.Context, id uuid.UUID) (Product, error) {
	const q = `
		SELECT id, name, price, created_at, updated_at
		FROM products
		WHERE id = $1`

	var p Product
	err := d.db.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.Nama, &p.Harga, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, repository.ErrNotFound
		}
		return Product{}, fmt.Errorf("find product by id: %w", err)
	}

	return p, nil
}

func (d *DBProductRepository) Create(ctx context.Context, in Product) (Product, error) {
	const q = `
		INSERT INTO products (name, price)
		VALUES ($1, $2)
		RETURNING id, name, price, created_at, updated_at`

	var p Product
	err := d.db.QueryRow(ctx, q, in.Nama, in.Harga).Scan(
		&p.ID, &p.Nama, &p.Harga, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return Product{}, fmt.Errorf("create product: %w", err)
	}

	return p, nil
}

func (d *DBProductRepository) Update(ctx context.Context, id uuid.UUID, in Product) (Product, error) {
	const q = `
		UPDATE products
		SET name = $2, price = $3, updated_at = now()
		WHERE id = $1
		RETURNING id, name, price, created_at, updated_at`

	var p Product
	err := d.db.QueryRow(ctx, q, id, in.Nama, in.Harga).Scan(
		&p.ID, &p.Nama, &p.Harga, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Product{}, repository.ErrNotFound
		}
		return Product{}, fmt.Errorf("update product: %w", err)
	}

	return p, nil
}

func (d *DBProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM products WHERE id = $1`

	tag, err := d.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
