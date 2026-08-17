package user

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (d *DBUserRepository) FindAll(ctx context.Context) ([]User, error) {
	const q = `
		select id, email, name, role, created_at, updated_at
		FROM users
		ORDER BY id DESC
		LIMIT 100`
	rows, err := d.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[User])
	if err != nil {
		return nil, fmt.Errorf("collect users: %w", err)
	}

	return users, nil
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
