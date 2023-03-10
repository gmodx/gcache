package memcache

import (
	"sync"
	"time"

	"github.com/gmodx/gcache/abstract"
)

var _ abstract.ICache[interface{}] = (*MemCache[interface{}])(nil)

type MemCache[T interface{}] struct {
	mu       sync.RWMutex
	options  MemCacheOptions
	cacheMap map[string]*CacheEntity[T]
}

type MemCacheOptions struct {
	CleanupInterval time.Duration
}

var defaultOptions = MemCacheOptions{
	CleanupInterval: time.Minute,
}

// Get retrieves an item from the cache by its key.
func (mc *MemCache[T]) Get(key string) *T {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if val, ok := mc.cacheMap[key]; ok && !val.Expired() {
		return &val.item
	}

	return nil
}

// Refresh updates the expiration time of an item in the cache.
func (mc *MemCache[T]) Refresh(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if val, ok := mc.cacheMap[key]; ok && !val.Expired() {
		mc.cacheMap[key].expiredTime = time.Now().Add(val.expiration)
	}
}

// Set adds or updates an item in the cache by its key.
func (mc *MemCache[T]) Set(key string, val T, opts abstract.CacheEntryOptions) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	mc.cacheMap[key] = &CacheEntity[T]{
		expiredTime: now.Add(opts.Expiration),
		expiration:  opts.Expiration,
		item:        val,
	}
}

// Remove deletes an item from the cache by its key, if it exists.
func (mc *MemCache[T]) Remove(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.cacheMap, key)
}

// New creates a new MemCache with the given options.
func New[T interface{}](opts MemCacheOptions) *MemCache[T] {
	gc := &MemCache[T]{
		options:  opts,
		cacheMap: map[string]*CacheEntity[T]{},
	}

	go gc.startClearJob()
	return gc
}

// NewDefault creates a new MemCache with default options.
func NewDefault[T interface{}]() *MemCache[T] {
	return New[T](defaultOptions)
}

// startClearJob starts a goroutine that periodically deletes expired items from the cache.
func (mc *MemCache[T]) startClearJob() {
	t := time.NewTicker(mc.options.CleanupInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			mc.DeleteAllExpired()
		}
	}
}

// Flush removes all items from the cache.
func (mc *MemCache[T]) Flush() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cacheMap = map[string]*CacheEntity[T]{}
}

// DeleteAllExpired deletes all expired items from the cache.
func (gc *MemCache[T]) DeleteAllExpired() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	for cacheKey, cacheItem := range gc.cacheMap {
		if cacheItem.Expired() {
			delete(gc.cacheMap, cacheKey)
		}
	}
}
