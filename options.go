package memcache

import "time"

type Options struct {
	CleanupInterval time.Duration
	Expiration      time.Duration
}

var defaultOptions = Options{
	CleanupInterval: time.Minute,
	Expiration:      time.Minute,
}
