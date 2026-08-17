# Golang REST API

REST API dengan Fiber v3, PostgreSQL 18, dan pgx.

## Prasyarat

| Tool | Versi | Cek |
|---|---|---|
| Go | 1.25+ | `go version` |
| Docker | dengan Compose v2 | `docker compose version` |
| Make | — | `make --version` |

Go dijalankan di host, hanya PostgreSQL yang di dalam Docker.

## Setup Pertama Kali

```bash
git clone <url-repo>
cd golang-rest-api

make setup
```

Perintah tersebut akan:
1. Membuat `.env` dari `.env.example`
2. Download dependency Go + install `migrate` dan `air`
3. Menyalakan container PostgreSQL
4. Menunggu healthcheck sampai `healthy`
5. Menjalankan migration

Kalau `migrate` atau `air` baru terinstall, pastikan `$(go env GOPATH)/bin` ada di `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Tambahkan ke `~/.zshrc` agar permanen.

## Menjalankan

```bash
make dev     # dengan hot reload (air)
make run     # tanpa hot reload
```

Server: http://localhost:3000

## Perintah Lain

| Perintah | Fungsi |
|---|---|
| `make db-up` / `db-down` | nyalakan / matikan Postgres |
| `make db-shell` | masuk psql |
| `make db-logs` | pantau log Postgres |
| `make db-reset` | ⚠️ hapus volume + migrate ulang |
| `make migrate-up` | jalankan migration |
| `make migrate-down` | mundur 1 migration |
| `make migrate-create name=xxx` | buat file migration baru |
| `make test` | jalankan test |
| `make fmt` | format + vet |

## Endpoint

```
GET  /api/health          liveness
GET  /api/ready           readiness + status pool
GET  /api/users
GET  /api/users/:id
POST /api/users
GET  /api/products
GET  /api/products/:id
POST /api/products
```

## Catatan

**Port 5433, bukan 5432.** Compose memetakan `127.0.0.1:5433:5432` agar tidak bentrok dengan Postgres lain di host. Binding ke `127.0.0.1` berarti database hanya bisa diakses dari mesin ini.

**UUIDv7 digenerate PostgreSQL** lewat `DEFAULT uuidv7()` (fitur PostgreSQL 18). Aplikasi menerimanya via `RETURNING`.

## Troubleshooting

**`connection refused` di port 5433**
```bash
docker compose ps    # kolom PORTS harus 127.0.0.1:5433->5432/tcp
```

**`Dirty database version N`** — migration gagal di tengah. Periksa kondisi schema, lalu:
```bash
migrate -path db/migrations -database "$DB_URL" force <versi-terakhir-yang-benar>
```

**Data hilang setelah `docker compose down`** — pastikan volume ter-mount ke `/var/lib/postgresql` (path PostgreSQL 18), bukan `/var/lib/postgresql/data`.