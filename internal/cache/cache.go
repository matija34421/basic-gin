package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (data []byte, ok bool, err error)
	Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
	Del(ctx context.Context, keys ...string) error
}
