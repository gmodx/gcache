package abstract

import "time"

// ICache is an interface that defines the basic operations of a cache.
type ICache[T interface{}] interface {
	Get(key string) *T
	Set(key string, val T, opts CacheEntryOptions)
	Remove(key string)
	Refresh(key string)
}

// CacheEntryOptions represents options for a cached item.
type CacheEntryOptions struct {
	Expiration time.Duration // Expiration specifies the duration until the cached item expires.
}

const (
	NoExpiration = 100 * 365 * 24 * time.Hour // NoExpiration represents a duration indicating that a cached item has no expiration.
)
