package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"golang-rest-api/internal/repository"
)

// fakeUserRepo memenuhi interface UserRepository tanpa database
type fakeUserRepo struct {
	// Perilaku yang bisa diatur per test
	createFn   func(ctx context.Context, u User) (User, error)
	findByIDFn func(ctx context.Context, id uuid.UUID) (User, error)
	findAllFn  func(ctx context.Context) ([]User, error)
	updateFn   func(ctx context.Context, id uuid.UUID, u User) (User, error)
	deleteFn   func(ctx context.Context, id uuid.UUID) error
	// Perekam untuk verifikasi
	createCalled bool
	createInput  User
}

func (f *fakeUserRepo) Create(ctx context.Context, u User) (User, error) {
	f.createCalled = true
	f.createInput = u
	if f.createFn != nil {
		return f.createFn(ctx, u)
	}
	u.ID = uuid.New()
	return u, nil
}

func (f *fakeUserRepo) FindById(ctx context.Context, id uuid.UUID) (User, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return User{}, repository.ErrNotFound
}

func (f *fakeUserRepo) FindAll(ctx context.Context) ([]User, error) {
	if f.findAllFn != nil {
		return f.findAllFn(ctx)
	}
	return nil, nil
}

func (f *fakeUserRepo) Update(ctx context.Context, id uuid.UUID, u User) (User, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, u)
	}
	u.ID = id
	return u, nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func TestCreateUser_Validasi(t *testing.T) {
	tests := []struct {
		nama     string
		input    User
		wantErr  error
		wantCall bool // apakah repo seharusnya dipanggil?
	}{
		{
			nama:     "email kosong ditolak",
			input:    User{Nama: "Budi", Email: ""},
			wantErr:  repository.ErrValidation,
			wantCall: false,
		},
		{
			nama:     "nama kosong ditolak",
			input:    User{Nama: "", Email: "budi@mail.com"},
			wantErr:  repository.ErrValidation,
			wantCall: false,
		},
		{
			nama:     "email hanya spasi ditolak",
			input:    User{Nama: "Budi", Email: "   "},
			wantErr:  repository.ErrValidation,
			wantCall: false,
		},
		{
			nama:     "input valid diteruskan",
			input:    User{Nama: "Budi", Email: "budi@mail.com"},
			wantErr:  nil,
			wantCall: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			repo := &fakeUserRepo{}
			uc := NewUserUseCase(repo)

			_, err := uc.CreateUser(context.Background(), tc.input)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("mau error %v, dapat %v", tc.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("tidak mau error, dapat %v", err)
			}

			if repo.createCalled != tc.wantCall {
				t.Errorf("createCalled = %v, mau %v", repo.createCalled, tc.wantCall)
			}
		})
	}
}

func TestCreateUser_TrimSpasi(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUserUseCase(repo)

	_, err := uc.CreateUser(context.Background(), User{
		Nama:  "  Budi  ",
		Email: "  budi@mail.com  ",
	})
	if err != nil {
		t.Fatalf("tidak mau error: %v", err)
	}

	if repo.createInput.Nama != "Budi" {
		t.Errorf("nama = %q, mau %q", repo.createInput.Nama, "Budi")
	}
	if repo.createInput.Email != "budi@mail.com" {
		t.Errorf("email = %q, mau %q", repo.createInput.Email, "budi@mail.com")
	}
}
func TestCreateUser_ErrorDariRepo(t *testing.T) {
	repo := &fakeUserRepo{
		createFn: func(ctx context.Context, u User) (User, error) {
			return User{}, repository.ErrDuplicateKey
		},
	}
	uc := NewUserUseCase(repo)

	_, err := uc.CreateUser(context.Background(), User{
		Nama:  "Budi",
		Email: "budi@mail.com",
	})

	if !errors.Is(err, repository.ErrDuplicateKey) {
		t.Fatalf("mau ErrDuplicateKey, dapat %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	repo := &fakeUserRepo{
		deleteFn: func(ctx context.Context, id uuid.UUID) error {
			return repository.ErrNotFound
		},
	}
	uc := NewUserUseCase(repo)

	err := uc.DeleteUser(context.Background(), uuid.New())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("mau ErrNotFound, dapat %v", err)
	}
}
