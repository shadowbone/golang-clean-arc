package config

import (
	"net"
	"net/url"
	"time"
)

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	AppName         string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	ConnectTimeout  time.Duration
	TraceQuery      bool
}

func (d DBConfig) DSN() string {
	q := url.Values{}
	q.Set("sslmode", d.SSLMode)

	if d.AppName != "" {
		q.Set("application_name", d.AppName)
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(d.User, d.Password),
		Host:     net.JoinHostPort(d.Host, d.Port),
		Path:     "/" + d.Name,
		RawQuery: q.Encode(),
	}

	return u.String()
}
