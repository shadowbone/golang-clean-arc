package user

import (
	"context"
	"golang-rest-api/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (u *UserUseCase) GetAllUser(ctx context.Context) ([]User, error) {
	return u.repo.FindAll(ctx)
}

func (u *UserUseCase) FindByIdUser(ctx context.Context, id uuid.UUID) (User, error) {
	return u.repo.FindById(ctx, id)
}

func (u *UserUseCase) CreateUser(ctx context.Context, user User) (User, error) {
	user.Email = strings.TrimSpace(user.Email)
	user.Nama = strings.TrimSpace(user.Nama)

	if user.Email == "" {
		return User{}, repository.Validationf("email tidak boleh kosong")
	}

	if user.Nama == "" {
		return User{}, repository.Validationf("nama tidak boleh kosong")
	}

	return u.repo.Create(ctx, user)
}

func (u *UserUseCase) UpdateUser(ctx context.Context, id uuid.UUID, in User) (User, error) {
	in.Email = strings.TrimSpace(in.Email)
	in.Nama = strings.TrimSpace(in.Nama)

	if in.Email == "" {
		return User{}, repository.Validationf("email tidak boleh kosong")
	}

	if in.Nama == "" {
		return User{}, repository.Validationf("nama tidak boleh kosong")
	}

	return u.repo.Update(ctx, id, in)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
