package memcache

type IMemCache[T interface{}] interface {
	Get(key string) *T
	Set(key string, val T)
	Remove(key string)
	Refresh(key string)
}
