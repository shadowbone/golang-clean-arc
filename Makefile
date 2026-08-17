include .env
export

DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=$(DB_SSL)
MIGRATE = migrate -path db/migrations -database "$(DB_URL)"

## ============ SETUP ============

## setup pertama kali: env, dependency, database, migration
setup: env deps db-up wait-db migrate-up
	@echo ""
	@echo "✅ Setup selesai. Jalankan: make dev"

## buat .env dari .env.example kalau belum ada
env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ .env dibuat dari .env.example"; \
	else \
		echo "⏭  .env sudah ada, dilewati"; \
	fi

## download dependency Go + install CLI tools
deps:
	go mod download
	@command -v migrate >/dev/null 2>&1 || { \
		echo "📦 Install golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	}
	@command -v air >/dev/null 2>&1 || { \
		echo "📦 Install air..."; \
		go install github.com/air-verse/air@latest; \
	}
	@echo "✅ Dependency siap"

## ============ DATABASE ============

## nyalakan container postgres
db-up:
	docker compose up -d

## matikan container (data tetap aman)
db-down:
	docker compose down

## ⚠️ hapus volume + migrate ulang, SEMUA DATA HILANG
db-reset:
	docker compose down -v
	docker compose up -d
	@$(MAKE) wait-db
	@$(MAKE) migrate-up

## pantau log postgres
db-logs:
	docker compose logs -f postgres

## masuk ke psql
db-shell:
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

## tunggu postgres siap menerima koneksi
wait-db:
	@echo "⏳ Menunggu Postgres siap..."
	@for i in $$(seq 1 30); do \
		if docker compose exec -T postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) >/dev/null 2>&1; then \
			echo "✅ Postgres siap"; exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "❌ Timeout menunggu Postgres"; exit 1

## ============ MIGRATION ============

## ini untuk menjalankan up migration
migrate-up:
	$(MIGRATE) up

## ini untuk menjalankan down/membalikan migration
migrate-down:
	$(MIGRATE) down 1

## cek versi migration saat ini
migrate-status:
	$(MIGRATE) version

## buat file migration baru: make migrate-create name=create_orders_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Wajib isi name. Contoh: make migrate-create name=create_orders_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)

## ⚠️ paksa set versi saat database dirty: make migrate-force v=1
migrate-force:
	@if [ -z "$(v)" ]; then \
		echo "❌ Wajib isi v. Contoh: make migrate-force v=1"; \
		exit 1; \
	fi
	$(MIGRATE) force $(v)

## ============ RUN ============

## jalankan dengan hot reload
dev:
	air

## jalankan tanpa hot reload
run:
	go run ./cmd/web

## compile ke bin/app
build:
	go build -o bin/app ./cmd/web

## ============ QUALITY ============

test:
	go test ./... -v

test-cover:
	go test ./... -coverprofile=cover.out
	@go tool cover -func=cover.out | tail -1

## format + static check
fmt:
	go fmt ./...
	go vet ./...

tidy:
	go mod tidy

## bersihkan hasil build
clean:
	rm -rf bin tmp cover.out

## tampilkan daftar perintah
help:
	@grep -B1 -E '^[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) \
		| grep -A1 '^##' \
		| sed 's/^## //' \
		| paste - - -d'|' \
		| awk -F'|' '{printf "  \033[36m%-18s\033[0m %s\n", $$2, $$1}' \
		| sed 's/://'

## fungsinya untuk memberitahu make bahwa target yang terdaftar bukan nama file
.PHONY: setup env deps db-up db-down db-reset db-logs db-shell wait-db \
        migrate-up migrate-down migrate-status migrate-create migrate-force \
        dev run build test test-cover fmt tidy clean help