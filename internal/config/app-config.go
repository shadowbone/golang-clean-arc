package config

import (
	"strings"
	"time"
)

type AppConfig struct {
	Port            string
	Env             string
	ShutdownTimeout time.Duration
}

func (a AppConfig) IsProduction() bool {
	return strings.EqualFold(a.Env, "production")
}
