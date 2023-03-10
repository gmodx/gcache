package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/redis/go-redis/v9"
)

var _ abstract.ICache[interface{}] = (*RedisCache[interface{}])(nil)

type RedisCache[T interface{}] struct {
	options  Options
	dbClient *redis.Client
}

type CacheEntity[T interface{}] struct {
	expiredTime time.Time
	expiration  time.Duration
	item        T
}

// Get implements abstract.ICache
func (c *RedisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	val, err := c.dbClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}

// Refresh implements abstract.ICache
func (*RedisCache[T]) Refresh(ctx context.Context, key string) error {
	panic("unimplemented")
}

// Remove implements abstract.ICache
func (*RedisCache[T]) Remove(ctx context.Context, key string) error {
	panic("unimplemented")
}

// Set implements abstract.ICache
func (c *RedisCache[T]) Set(ctx context.Context, key string, val T, opts abstract.CacheEntryOptions) error {
	expiration := opts.Expiration
	if opts.Expiration == abstract.NoExpiration {
		expiration = redis.KeepTTL
	}

	return c.dbClient.Set(ctx, key, val, expiration).Err()
}

func (c *RedisCache[T]) initClient() {
	c.dbClient = redis.NewClient(&c.options.ConnectionOptions)
}

type Options struct {
	CleanupInterval   time.Duration
	ConnectionOptions redis.Options
}

var defaultOptions = Options{
	CleanupInterval: time.Minute,
	ConnectionOptions: redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	},
}

func New[T interface{}](opts Options) *RedisCache[T] {
	gc := &RedisCache[T]{
		options: opts,
	}

	gc.initClient()
	return gc
}

func NewDefault[T interface{}]() *RedisCache[T] {
	return New[T](defaultOptions)
}
