package address

import (
	"context"
	"golang-rest-api/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type UseCase struct {
	repo AddressRepository
	tx   repository.Transactor
}

func NewUseCase(
	repo AddressRepository,
	tx repository.Transactor,
) *UseCase {
	return &UseCase{repo: repo, tx: tx}
}

// Create
func (u *UseCase) Create(ctx context.Context, userId uuid.UUID, in AddressRequest) (Address, error) {
	in.Label = strings.TrimSpace(in.Label)
	in.Street = strings.TrimSpace(in.Street)
	in.City = strings.TrimSpace(in.City)

	if in.Label == "" {
		return Address{}, repository.Validationf("Label tidak boleh kosong")
	}

	if in.Street == "" {
		return Address{}, repository.Validationf("Street tidak boleh kosong")
	}

	if in.City == "" {
		return Address{}, repository.Validationf("City tidak boleh kosong")
	}

	input := in.ToEntity()
	input.UserID = userId

	var created Address
	err := u.tx.WithinTx(ctx, func(ctx context.Context) error {
		count, err := u.repo.CountByUserID(ctx, userId)
		if err != nil {
			return err
		}

		input.IsDefault = count == 0
		created, err = u.repo.Create(ctx, input)
		return err
	})

	return created, err
}

// FindByUserID
func (u *UseCase) FindByIdUser(ctx context.Context, userId uuid.UUID) ([]Address, error) {
	return u.repo.FindByUserID(ctx, userId)
}

func (u *UseCase) Update(ctx context.Context, userId, id uuid.UUID, in AddressRequest) (Address, error) {
	in.Label = strings.TrimSpace(in.Label)
	in.Street = strings.TrimSpace(in.Street)
	in.City = strings.TrimSpace(in.City)

	if in.Label == "" {
		return Address{}, repository.Validationf("Label tidak boleh kosong")
	}

	if in.Street == "" {
		return Address{}, repository.Validationf("Street tidak boleh kosong")
	}

	if in.City == "" {
		return Address{}, repository.Validationf("City tidak boleh kosong")
	}

	input := in.ToEntity()
	return u.repo.Update(ctx, userId, id, input)
}

func (u *UseCase) Delete(ctx context.Context, userId, id uuid.UUID) error {
	return u.repo.Delete(ctx, userId, id)
}

func (u *UseCase) SetDefault(ctx context.Context, userId, id uuid.UUID) error {
	return u.tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.repo.UnsetDefault(ctx, userId); err != nil {
			return err
		}

		return u.repo.SetDefault(ctx, userId, id)
	})
}
