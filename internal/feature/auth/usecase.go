package auth

import (
	"context"
	"errors"
	"golang-rest-api/internal/logger"
	"golang-rest-api/internal/repository"
	"golang-rest-api/internal/validation"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

type AuthUseCase struct {
	repo   AuthRepository
	users  UserProvider
	tokens *TokenManager
	hasher *Hasher
	tx     repository.Transactor
}

func NewAuthUseCase(
	repo AuthRepository,
	users UserProvider,
	tokens *TokenManager,
	hasher *Hasher,
	tx repository.Transactor,
) *AuthUseCase {
	return &AuthUseCase{repo: repo, users: users, tokens: tokens, hasher: hasher, tx: tx}
}

func (u *AuthUseCase) Register(ctx context.Context, req RegisterRequest) (TokenResponse, error) {
	if err := validation.Struct(req); err != nil {
		return TokenResponse{}, err
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Nama = strings.TrimSpace(req.Nama)

	hash, err := u.hasher.Hash(req.Password)
	if err != nil {
		return TokenResponse{}, err
	}

	var resp TokenResponse

	err = u.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err := u.users.Create(ctx, req.Nama, req.Email)
		if err != nil {
			return err
		}

		if err := u.repo.CreateCredential(ctx, created.ID, hash); err != nil {
			return err
		}

		resp, err = u.issueTokens(ctx, created.ID)
		return err
	})

	return resp, err
}

func (u *AuthUseCase) Login(ctx context.Context, req LoginRequest) (TokenResponse, error) {
	if err := validation.Struct(req); err != nil {
		return TokenResponse{}, err
	}

	info, err := u.users.FindByEmail(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			u.hasher.VerifyDummy(req.Password)

			return TokenResponse{}, ErrInvalidCredentials
		}
		return TokenResponse{}, err
	}

	cred, err := u.repo.FindCredentialByUserID(ctx, info.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			u.hasher.VerifyDummy(req.Password)
			return TokenResponse{}, ErrInvalidCredentials
		}
		return TokenResponse{}, err
	}

	if !u.hasher.Verify(cred.PasswordHash, req.Password) {
		return TokenResponse{}, ErrInvalidCredentials
	}

	return u.issueTokens(ctx, info.ID)
}

func (u *AuthUseCase) Refresh(ctx context.Context, req RefreshRequest) (TokenResponse, error) {
	if err := validation.Struct(req); err != nil {
		return TokenResponse{}, err
	}

	log := logger.FromContext(ctx)
	hash := HasToken(req.RefreshToken)

	var resp TokenResponse

	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		old, err := u.repo.FindRefreshTokenByHash(ctx, hash)
		if err != nil {
			return err
		}

		if old.RevokedAt != nil {
			log.Warn("refresh token reuse terdeteksi",
				slog.String("user_id", old.UserID.String()),
				slog.String("token_id", old.ID.String()),
			)

			if err := u.repo.RevokeAllByUserID(ctx, old.UserID); err != nil {
				return err
			}
			return ErrTokenRevoked
		}

		if !old.IsUsable() {
			return ErrTokenExpired
		}

		resp, err = u.issueTokens(ctx, old.UserID)
		if err != nil {
			return err
		}

		newHash := HasToken(resp.RefreshToken)
		created, err := u.repo.FindRefreshTokenByHash(ctx, newHash)
		if err != nil {
			return err
		}

		return u.repo.RevokeToken(ctx, old.ID, &created.ID)
	})

	return resp, err
}

func (u *AuthUseCase) Logout(ctx context.Context, req RefreshRequest) error {
	if err := validation.Struct(req); err != nil {
		return err
	}

	t, err := u.repo.FindRefreshTokenByHash(ctx, HasToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return nil
		}
		return err
	}

	return u.repo.RevokeToken(ctx, t.ID, nil)
}

func (u *AuthUseCase) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return u.repo.RevokeAllByUserID(ctx, userID)
}

func (u *AuthUseCase) issueTokens(ctx context.Context, userID uuid.UUID) (TokenResponse, error) {
	access, exp, err := u.tokens.GenerateAccess(userID)
	if err != nil {
		return TokenResponse{}, err
	}

	raw, hash, refreshExp, err := u.tokens.GenerateRefresh()
	if err != nil {
		return TokenResponse{}, err
	}

	_, err = u.repo.SaveRefreshToken(ctx, RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: refreshExp,
	})

	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresAt:    exp,
		TokenType:    "Bearer",
	}, nil
}
