package auth

import (
	"context"
	"errors"
	"fmt"
	"golang-rest-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	CreateCredential(ctx context.Context, userID uuid.UUID, hash string) error
	FindCredentialByUserID(ctx context.Context, userID uuid.UUID) (UserCredential, error)

	SaveRefreshToken(ctx context.Context, t RefreshToken) (RefreshToken, error)
	FindRefreshTokenByHash(ctx context.Context, hash string) (RefreshToken, error)
	RevokeToken(ctx context.Context, id uuid.UUID, replacedBy *uuid.UUID) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error
	RevokeChain(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type DBAuthRepository struct {
	repository.DB
}

func NewAuthRepository(pool repository.Querier) AuthRepository {
	return &DBAuthRepository{DB: repository.NewDB(pool)}
}

func (d *DBAuthRepository) CreateCredential(ctx context.Context, userID uuid.UUID, hash string) error {
	const q = `
		INSERT INTO user_credentials (user_id, password_hash)
		VALUES ($1, $2)
	`

	if _, err := d.Q(ctx).Exec(ctx, q, userID, hash); err != nil {
		return fmt.Errorf("create credential: %w", err)
	}
	return nil
}

func (d *DBAuthRepository) FindCredentialByUserID(ctx context.Context, userID uuid.UUID) (UserCredential, error) {
	const q = `
		SELECT user_id, password_hash, password_changed_at
		FROM user_credentials
		WHERE user_id = $1
	`

	var uc UserCredential
	err := d.Q(ctx).QueryRow(ctx, q, userID).Scan(
		&uc.UserID, &uc.PasswordHash, &uc.PasswordChangedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserCredential{}, repository.ErrNotFound
		}
		return UserCredential{}, fmt.Errorf("find credential: %w", err)
	}

	return uc, nil
}

func (d *DBAuthRepository) SaveRefreshToken(ctx context.Context, t RefreshToken) (RefreshToken, error) {
	const q = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, revoked_at, replaced_by, created_at`

	var out RefreshToken
	err := d.Q(ctx).QueryRow(ctx, q, t.UserID, t.TokenHash, t.ExpiresAt).Scan(
		&out.ID, &out.UserID, &out.TokenHash, &out.ExpiresAt,
		&out.RevokedAt, &out.ReplacedBy, &out.CreatedAt,
	)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("save refresh token: %w", err)
	}

	return out, nil
}

func (d *DBAuthRepository) FindRefreshTokenByHash(ctx context.Context, hash string) (RefreshToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, revoked_at, replaced_by, created_at
		FROM refresh_tokens
		WHERE token_hash = $1`

	var t RefreshToken
	err := d.Q(ctx).QueryRow(ctx, q, hash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt,
		&t.RevokedAt, &t.ReplacedBy, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshToken{}, ErrInvalidToken
		}
		return RefreshToken{}, fmt.Errorf("find refresh token: %w", err)
	}

	return t, nil
}

func (d *DBAuthRepository) RevokeToken(ctx context.Context, id uuid.UUID, replacedBy *uuid.UUID) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = now(), replaced_by = $2
		WHERE id = $1 AND revoked_at IS NULL`

	if _, err := d.Q(ctx).Exec(ctx, q, id, replacedBy); err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

func (d *DBAuthRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID) error {
	const q = `
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE user_id = $1 AND revoked_at IS NULL`

	if _, err := d.Q(ctx).Exec(ctx, q, userID); err != nil {
		return fmt.Errorf("revoke all tokens: %w", err)
	}
	return nil
}

func (d *DBAuthRepository) RevokeChain(ctx context.Context, tokenID uuid.UUID) error {
	const q = `
		WITH RECURSIVE chain AS (
			SELECT id, replaced_by FROM refresh_tokens WHERE id = $1
			UNION ALL
			SELECT rt.id, rt.replaced_by
			FROM refresh_tokens rt
			JOIN chain c ON rt.id = c.replaced_by
		)
		UPDATE refresh_tokens
		SET revoked_at = now()
		WHERE id IN (SELECT id FROM chain) AND revoked_at IS NULL`

	if _, err := d.Q(ctx).Exec(ctx, q, tokenID); err != nil {
		return fmt.Errorf("revoke chain: %w", err)
	}
	return nil
}

func (d *DBAuthRepository) DeleteExpired(ctx context.Context) (int64, error) {
	const q = `
		DELETE FROM refresh_tokens
		WHERE expires_at < now() - interval '30 days'`

	tag, err := d.Q(ctx).Exec(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("delete expired tokens: %w", err)
	}
	return tag.RowsAffected(), nil
}
