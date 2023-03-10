package abstract

import (
	"context"
	"time"
)

// ICache is an interface that defines the basic operations of a cache.
type ICache[T interface{}] interface {
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, val T, opts CacheEntryOptions) error
	Remove(ctx context.Context, key string) error
	Refresh(ctx context.Context, key string) error
}

// CacheEntryOptions represents options for a cached item.
type CacheEntryOptions struct {
	Expiration time.Duration // Expiration specifies the duration until the cached item expires.
}

const (
	NoExpiration = 100 * 365 * 24 * time.Hour // NoExpiration represents a duration indicating that a cached item has no expiration.
)
