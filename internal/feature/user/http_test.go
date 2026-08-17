package user

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"golang-rest-api/internal/repository"
	"golang-rest-api/internal/response"
)

func setupApp(repo UserRepository) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: response.ErrorHandler,
	})

	h := NewUserHandler(NewUserUseCase(repo))

	api := app.Group("/api/users")
	api.Get("/", h.FetchAll)
	api.Get("/:id", h.FetchById)
	api.Post("/", h.Store)

	return app
}

func TestStore_StatusCode(t *testing.T) {
	tests := []struct {
		nama     string
		body     string
		repo     *fakeUserRepo
		wantCode int
	}{
		{
			nama:     "sukses",
			body:     `{"nama":"Budi","email":"budi@mail.com"}`,
			repo:     &fakeUserRepo{},
			wantCode: http.StatusCreated,
		},
		{
			nama:     "validasi gagal",
			body:     `{"nama":"","email":""}`,
			repo:     &fakeUserRepo{},
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			nama: "email duplikat",
			body: `{"nama":"Budi","email":"budi@mail.com"}`,
			repo: &fakeUserRepo{
				createFn: func(ctx context.Context, u User) (User, error) {
					return User{}, repository.ErrDuplicateKey
				},
			},
			wantCode: http.StatusConflict,
		},
		{
			nama:     "JSON rusak",
			body:     `{invalid`,
			repo:     &fakeUserRepo{},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.nama, func(t *testing.T) {
			app := setupApp(tc.repo)

			req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantCode {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("status = %d, mau %d. body: %s", resp.StatusCode, tc.wantCode, body)
			}
		})
	}
}

func TestFetchById_UUIDInvalid(t *testing.T) {
	app := setupApp(&fakeUserRepo{})

	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/api/users/bukan-uuid", nil))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, mau 400", resp.StatusCode)
	}
}

func TestFetchById_NotFound(t *testing.T) {
	app := setupApp(&fakeUserRepo{}) // default FindById → ErrNotFound

	url := "/api/users/" + uuid.New().String()
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, url, nil))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, mau 404", resp.StatusCode)
	}

	var body map[string]any
	json.NewDecoder(resp.Body).Decode(&body)

	if body["status"] != "gagal" {
		t.Errorf(`status field = %v, mau "gagal"`, body["status"])
	}
}
