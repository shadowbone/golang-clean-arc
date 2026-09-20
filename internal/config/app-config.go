package config

import (
	"strings"
	"time"
)

type AppConfig struct {
	Port            string
	Env             string
	ShutdownTimeout time.Duration
	TrustedProxies  []string
	MetricsEnabled  bool
	MetricsPath     string
}

func (a AppConfig) IsProduction() bool {
	return strings.EqualFold(a.Env, "production")
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}
