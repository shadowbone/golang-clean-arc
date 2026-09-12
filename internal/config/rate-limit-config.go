package config

import "time"

type RateLimitingConfig struct {
	Enabled        bool
	LoginMax       int
	LoginWindow    time.Duration
	RegisterMax    int
	RegisterWindow time.Duration
}
