package auth

import (
	"context"

	"golang-rest-api/internal/feature/user"
)

type userAdapter struct {
	repo user.UserRepository
}

func NewUserAdapter(repo user.UserRepository) UserProvider {
	return &userAdapter{repo: repo}
}

func (a *userAdapter) FindByEmail(ctx context.Context, email string) (UserInfo, error) {
	u, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{ID: u.ID, Nama: u.Nama, Email: u.Email}, nil
}

func (a *userAdapter) Create(ctx context.Context, nama, email string) (UserInfo, error) {
	u, err := a.repo.Create(ctx, user.User{Nama: nama, Email: email})
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{ID: u.ID, Nama: u.Nama, Email: u.Email}, nil
}
