package database

import (
	"context"
	"fmt"
	"golang-rest-api/internal/config"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

func NewPool(c context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = time.Minute
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	poolCfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}

	if cfg.TraceQuery {
		poolCfg.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   &slogAdapter{log: slog.Default()},
			LogLevel: tracelog.LogLevelDebug,
		}
	}

	pool, err := pgxpool.NewWithConfig(c, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("buat pool:%w", err)
	}

	pingCtx, cancel := context.WithTimeout(c, cfg.ConnectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database:%w", err)
	}

	return pool, nil
}

type slogAdapter struct {
	log *slog.Logger
}

func (a *slogAdapter) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]any, 0, len(data))
	for k, v := range data {
		if k == "sql" {
			if s, ok := v.(string); ok {
				v = strings.Join(strings.Fields(s), " ")
			}
		}
		attrs = append(attrs, slog.Any(k, v))
	}
	a.log.Debug(msg, attrs...)
}
