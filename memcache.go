package memcache

import (
	"sync"
	"time"
)

var _ IMemCache[interface{}] = (*MemCache[interface{}])(nil)

type MemCache[T interface{}] struct {
	mu       sync.RWMutex
	options  Options
	cacheMap map[string]*CacheEntity[T]
}

// Get implements IMemCache
func (mc *MemCache[T]) Get(key string) *T {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if val, ok := mc.cacheMap[key]; ok && !val.Expired() {
		return &val.item
	}

	return nil
}

func (mc *MemCache[T]) Refresh(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if val, ok := mc.cacheMap[key]; ok && !val.Expired() {
		now := time.Now()
		mc.cacheMap[key].expiredTime = now.Add(mc.options.Expiration)
	}
}

func (mc *MemCache[T]) Set(key string, val T) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	mc.cacheMap[key] = &CacheEntity[T]{
		expiredTime: now.Add(mc.options.Expiration),
		item:        val,
	}
}

// Delete one cache if existed
func (mc *MemCache[T]) Remove(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.cacheMap, key)
}

func New[T interface{}](opts Options) *MemCache[T] {
	gc := &MemCache[T]{
		options:  opts,
		cacheMap: map[string]*CacheEntity[T]{},
	}

	go gc.startClearJob()
	return gc
}

func NewDefault[T interface{}]() *MemCache[T] {
	return New[T](defaultOptions)
}

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

// Clear all cache
func (mc *MemCache[T]) Flush() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.cacheMap = map[string]*CacheEntity[T]{}
}

// Delete all expired items.
func (gc *MemCache[T]) DeleteAllExpired() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	for cacheKey, cacheItem := range gc.cacheMap {
		if cacheItem.Expired() {
			delete(gc.cacheMap, cacheKey)
		}
	}
}
