package product

import (
	"context"
	"golang-rest-api/internal/repository"
	"strings"

	"github.com/google/uuid"
)

type ProductUseCase struct {
	repo ProductRepository
}

func NewProductUseCase(repo ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (p *ProductUseCase) GetAllProduct(ctx context.Context) ([]Product, error) {
	return p.repo.FindAll(ctx)
}

func (p *ProductUseCase) FindByIdProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	return p.repo.FindById(ctx, id)
}

func (p *ProductUseCase) CreateProduct(ctx context.Context, in Product) (Product, error) {
	in.Nama = strings.TrimSpace(in.Nama)

	if in.Nama == "" {
		return Product{}, repository.Validationf("nama tidak boleh kosong")
	}
	if in.Harga <= 0 {
		return Product{}, repository.Validationf("harga harus lebih dari 0")
	}

	return p.repo.Create(ctx, in)
}
