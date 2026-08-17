package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

func NewPool(c context.Context) (*pgxpool.Pool, error) {
	dsn, err := BuildDSN()
	if err != nil {
		return nil, err
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = time.Minute
	cfg.ConnConfig.ConnectTimeout = 5 * time.Second
	cfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}
	pool, err := pgxpool.NewWithConfig(c, cfg)
	if err != nil {
		return nil, fmt.Errorf("buat pool:%w", err)
	}

	pingCtx, cancel := context.WithTimeout(c, 5*time.Second)

	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database:%w", err)
	}

	return pool, nil
}

func BuildDSN() (string, error) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5433")
	name := os.Getenv("POSTGRES_DB")
	ssl := getEnv("DB_SSL", "disable")

	for k, v := range map[string]string{
		"POSTGRES_USER":     user,
		"POSTGRES_PASSWORD": pass,
		"POSTGRES_DB":       name,
	} {
		if v == "" {
			return "", fmt.Errorf("env %s wajib diisi", k)
		}
	}

	q := url.Values{}
	q.Set("sslmode", ssl)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     net.JoinHostPort(host, port),
		Path:     "/" + name,
		RawQuery: q.Encode(),
	}

	return u.String(), nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
