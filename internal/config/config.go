package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	App       AppConfig
	DB        DBConfig
	Log       LogConfig
	Auth      AuthConfig
	Redis     RedisConfig
	RateLimit RateLimitingConfig
}

func Load() (*Config, error) {
	l := &loader{}

	cfg := &Config{
		App: AppConfig{
			Port:            l.str("APP_PORT", "3000"),
			Env:             l.str("APP_ENV", "development"),
			ShutdownTimeout: l.duration("SHUTDOWN_TIMEOUT", 20*time.Second),
		},
		DB: DBConfig{
			Host:            l.str("DB_HOST", "127.0.0.1"),
			Port:            l.str("DB_PORT", "5433"),
			User:            l.required("POSTGRES_USER"),
			Password:        l.required("POSTGRES_PASSWORD"),
			Name:            l.required("POSTGRES_DB"),
			SSLMode:         l.str("DB_SSL", "disable"),
			AppName:         l.str("DB_APP_NAME", "api"),
			MaxConns:        int32(l.integer("DB_MAX_CONNS", 10)),
			MinConns:        int32(l.integer("DB_MIN_CONNS", 2)),
			MaxConnLifetime: l.duration("DB_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime: l.duration("DB_MAX_CONN_IDLE", 30*time.Minute),
			ConnectTimeout:  l.duration("DB_CONNECT_TIMEOUT", 5*time.Second),
			TraceQuery:      l.boolean("DB_TRACE_QUERY", false),
		},
		Log: LogConfig{
			Level:       l.str("LOG_LEVEL", "info"),
			Format:      l.str("LOG_FORMAT", "json"),
			RequestBody: l.boolean("LOG_REQUEST_BODY", false),
		},
		Auth: AuthConfig{
			JWTSecret:       l.required("JWT_SECRET"),
			AccessTokenTTL:  l.duration("ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: l.duration("REFRESH_TOKEN_TTL", 15*24*time.Hour),
			BcryptCost:      l.integer("BYCRYPT_COST", 12),
		},
		Redis: RedisConfig{
			Addr:     l.str("REDIS_ADDR", "127.0.0.1:6380"),
			Password: l.str("REDIS_PASSWORD", ""),
			DB:       l.integer("REDIS_DB", 0),
		},
		RateLimit: RateLimitingConfig{
			Enabled:        l.boolean("RATE_LIMIT_ENABLED", true),
			LoginMax:       l.integer("RATE_LIMIT_LOGIN_MAX", 5),
			LoginWindow:    l.duration("RATE_LIMIT_LOGIN_WINDOW", 15*time.Minute),
			RegisterMax:    l.integer("RATE_LIMIT_REGISTER_MAX", 3),
			RegisterWindow: l.duration("RATE_LIMIT_REGISTER_WINDOW", time.Hour),
		},
	}
	if len(l.missing) > 0 {
		return nil, fmt.Errorf("env wajib belum diisi: %s", strings.Join(l.missing, ", "))
	}
	if len(l.errs) > 0 {
		return nil, fmt.Errorf("config tidak valid:\n  - %s", strings.Join(l.errs, "\n  - "))
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	switch strings.ToLower(c.Log.Level) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("LOG_LEVEL harus debug/info/warn/error, value %q", c.Log.Level)
	}

	switch strings.ToLower(c.Log.Format) {
	case "json", "text":
	default:
		return fmt.Errorf("LOG_FORMAT harus json/text, value %q", c.Log.Format)
	}

	if c.DB.MaxConns < c.DB.MinConns {
		return fmt.Errorf("DB_MAX_CONNS (%d) tidak boleh lebih kecil dari DB_MIN_CONNS (%d)",
			c.DB.MaxConns, c.DB.MinConns)
	}

	if c.App.IsProduction() {
		if c.DB.SSLMode == "disable" {
			return fmt.Errorf("DB_SSL tidak boleh 'disable' di production")
		}
		if c.Log.RequestBody {
			return fmt.Errorf("LOG_REQUEST_BODY tidak boleh aktif di production")
		}
	}

	if len(c.Auth.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET minimal 32 karakter")
	}

	if c.Auth.BcryptCost < 10 || c.Auth.BcryptCost > 15 {
		return fmt.Errorf("BCRYPT_COST harus diantara 10 dan 15")
	}

	return nil
}
