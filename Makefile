include .env
export

DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=$(DB_SSL)
MIGRATE = migrate -path db/migrations -database "$(DB_URL)"
GOLANGCI_VERSION = v2.13.1
IMAGE = golang-rest-api
VERSION ?= dev
## ============ SETUP ============


setup: env deps db-up wait-db migrate-up
	@echo ""
	@echo "Setup selesai. Jalankan: make dev"


env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env dibuat dari .env.example"; \
	else \
		echo "⏭  .env sudah ada, dilewati"; \
	fi


deps:
	go mod download
	@command -v migrate >/dev/null 2>&1 || { \
		echo "Install golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	}
	@command -v air >/dev/null 2>&1 || { \
		echo "Install air..."; \
		go install github.com/air-verse/air@latest; \
	}
	@echo "Dependency siap"

## ============ DATABASE ============


db-up:
	docker compose up -d


db-down:
	docker compose down


db-reset:
	docker compose down -v
	docker compose up -d
	@$(MAKE) wait-db
	@$(MAKE) migrate-up


db-logs:
	docker compose logs -f postgres


db-shell:
	docker compose exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)


wait-db:
	@echo "⏳ Menunggu Postgres siap..."
	@for i in $$(seq 1 30); do \
		if docker compose exec -T postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) >/dev/null 2>&1; then \
			echo "Postgres siap"; exit 0; \
		fi; \
		sleep 1; \
	done; \
	echo "❌ Timeout menunggu Postgres"; exit 1

## ============ MIGRATION ============

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-status:
	$(MIGRATE) version

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Wajib isi name. Contoh: make migrate-create name=create_orders_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-force:
	@if [ -z "$(v)" ]; then \
		echo "❌ Wajib isi v. Contoh: make migrate-force v=1"; \
		exit 1; \
	fi
	$(MIGRATE) force $(v)

## ============ RUN ============

dev:
	air

run:
	go run ./cmd/web


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


clean:
	rm -rf bin tmp cover.out

help:
	@grep -B1 -E '^[a-zA-Z0-9_-]+:' $(MAKEFILE_LIST) \
		| grep -A1 '^##' \
		| sed 's/^## //' \
		| paste - - -d'|' \
		| awk -F'|' '{printf "  \033[36m%-18s\033[0m %s\n", $$2, $$1}' \
		| sed 's/://'

lint:
	@command -v golangci-lint >/dev/null 2>&1 || \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	golangci-lint run

## jalankan semua pemeriksaan seperti di CI
check: fmt lint test
	@echo "✅ Semua pemeriksaan lolos"

## test dengan race detector, seperti di CI
test-race:
	go test -race ./...

## Build Docker
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$$(git rev-parse --short HEAD) \
		-t $(IMAGE):$(VERSION) .

## check ukuran dan user image
docker-inspect:
	@docker images $(IMAGE):$(VERSION)
	@docker run --rm --entrypoint sh $(IMAGE):$(VERSION) -c "id"

.PHONY: setup env deps db-up db-down db-reset db-logs db-shell wait-db \
        migrate-up migrate-down migrate-status migrate-create migrate-force \
        dev run build test test-cover fmt tidy clean help lint check test-race \
		docker-build docker-inspect