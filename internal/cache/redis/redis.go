package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr, pass string, db int) *Client {
	return &Client{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: pass,
			DB:       db,
		}),
	}
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, bool, error) {
	client, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return client, true, nil
}

func (c *Client) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) Ping(ctx context.Context) error { return c.rdb.Ping(ctx).Err() }

func (c *Client) Close() error { return c.rdb.Close() }
