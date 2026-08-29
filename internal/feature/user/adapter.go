package user

import (
	"context"
	"golang-rest-api/internal/feature/auth"
)

type AuthAdapter struct {
	repo UserRepository
}

func NewAuthAdapter(repo UserRepository) *AuthAdapter {
	return &AuthAdapter{repo: repo}
}

func (a *AuthAdapter) FindByEmail(ctx context.Context, email string) (auth.UserInfo, error) {
	u, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return auth.UserInfo{}, err
	}

	return auth.UserInfo{ID: u.ID, Nama: u.Nama, Email: u.Email}, nil
}

func (a *AuthAdapter) Create(ctx context.Context, nama, email string) (auth.UserInfo, error) {
	u, err := a.repo.Create(ctx, User{Nama: nama, Email: email})
	if err != nil {
		return auth.UserInfo{}, err
	}
	return auth.UserInfo{ID: u.ID, Nama: u.Nama, Email: u.Email}, nil
}
