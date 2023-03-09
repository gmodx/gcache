package memcache

import "time"

// CacheEntity represents a cached item in the cache.
type CacheEntity[T interface{}] struct {
	expiredTime time.Time
	expiration  time.Duration
	item        T
}

// Expired returns a boolean indicating whether the cached item has expired.
func (e *CacheEntity[T]) Expired() bool {
	return e.expiredTime.Before(time.Now())
}
