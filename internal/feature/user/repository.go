package user

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	repository.BaseRepository[User]
}

type DBUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &DBUserRepository{db: db}
}

func (d *DBUserRepository) FindAll(
	ctx context.Context,
	p repository.Pagination,
) (repository.Page[User], error) {
	const countQ = `SELECT COUNT(*) FROM users`

	var total int64
	if err := d.db.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return repository.Page[User]{}, fmt.Errorf("count user: %w", err)
	}
	const q = `
		select id, email, name, role, created_at, updated_at
		FROM users
		ORDER BY id DESC
		LIMIT $1 OFFSET $2`
	rows, err := d.db.Query(ctx, q, p.Limit, p.Offset())
	if err != nil {
		return repository.Page[User]{}, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[User])
	if err != nil {
		return repository.Page[User]{}, fmt.Errorf("collect users: %w", err)
	}

	return repository.NewPage(users, total, p), nil
}
func (d *DBUserRepository) FindById(ctx context.Context, id uuid.UUID) (User, error) {
	const q = `
		select id, email, name, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	var u User
	err := d.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.Nama, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, repository.ErrNotFound
		}
		return User{}, fmt.Errorf("find user by id: %w", err)
	}

	return u, nil
}

func (d *DBUserRepository) Create(ctx context.Context, in User) (User, error) {
	const q = `
		INSERT INTO users (email, name, role)
		VALUES ($1, $2, COALESCE(NULLIF($3, ''), 'user'))
		RETURNING id, email, name, role, created_at, updated_at`

	var u User
	err := d.db.QueryRow(ctx, q, in.Email, in.Nama, in.Role).Scan(
		&u.ID, &u.Email, &u.Nama, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (d *DBUserRepository) Update(ctx context.Context, id uuid.UUID, in User) (User, error) {
	const q = `
		UPDATE users
		SET name = $2, email = $3, role = $4, updated_at = now()
		WHERE id = $1
		RETURNING id, email, name, role, created_at, updated_at`

	var u User
	err := d.db.QueryRow(ctx, q, id, in.Nama, in.Email, in.Role).Scan(
		&u.ID, &u.Email, &u.Nama, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, repository.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, repository.ErrDuplicateKey
		}

		return User{}, fmt.Errorf("update user: %w", err)
	}

	return u, nil
}

func (d *DBUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM users WHERE id = $1`

	tag, err := d.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
