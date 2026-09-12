package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{client: client}
}

func (r *RedisLimiter) Allow(
	ctx context.Context,
	key string,
	max int,
	window time.Duration,
) (int, time.Duration, error) {
	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, window)
	ttl := pipe.TTL(ctx, key)

	if _, err := pipe.Exec(ctx); err != nil {
		return 0, 0, fmt.Errorf("rate limit check: %w", err)
	}

	count := int(incr.Val())
	remining := max - count

	return remining, ttl.Val(), nil
}
