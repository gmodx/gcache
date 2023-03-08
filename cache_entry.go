package memcache

import "time"

type CacheEntity[T interface{}] struct {
	expiredTime time.Time
	item        T
}

func (item *CacheEntity[T]) Expired() bool {
	return item.expiredTime.Before(time.Now())
}
