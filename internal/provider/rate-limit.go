package provider

import (
	"context"
	"fmt"
	"golang-rest-api/internal/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
)

type RateLimiter interface {
	Allow(ctx context.Context,
		key string,
		max int,
		window time.Duration,
	) (
		remaining int,
		retryAfter time.Duration,
		er error,
	)
}

type NoopLimiter struct{}

func (NoopLimiter) Allow(context.Context, string, int, time.Duration) (int, time.Duration, error) {
	return 1, 0, nil
}

func RateLimit(rl RateLimiter, prefix string, max int, window time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		key := fmt.Sprintf("rl:%s:%s", prefix, c.IP())

		remaining, retryAfter, err := rl.Allow(c.Context(), key, max, window)
		if err != nil {
			return c.Next()
		}

		c.Set("X-RateLimit-Limit", strconv.Itoa(max))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(maxInt(remaining, 0)))

		if remaining < 0 {
			c.Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
			return repository.ErrTooManyRequests
		}

		return c.Next()
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
