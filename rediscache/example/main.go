package main

import (
	"context"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/gmodx/gcache/rediscache"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.TODO()
	connOpts := redis.Options{
		Addr: "10.10.100.101:6379",
	}

	var c1 abstract.ICache[string] = rediscache.New[string](rediscache.Options{ConnectionOptions: connOpts})
	err := c1.Set(ctx, "baz", "bar", abstract.CacheEntryOptions{Expiration: abstract.NoExpiration})
	if err != nil {
		panic(err)
	}

	var c2 abstract.ICache[int] = rediscache.New[int](rediscache.Options{ConnectionOptions: connOpts})
	err = c2.Set(ctx, "baz2", 8, abstract.CacheEntryOptions{Expiration: 50 * time.Second})
	if err != nil {
		panic(err)
	}

	// // Get the value stored in the cache with the key "baz".
	// val := c.Get("baz")
	// if val != nil {
	// 	fmt.Println(val)
	// }

	// // Refresh the expiration time of the cached item with the key "baz".
	// c.Refresh("baz")

	// // Remove the cached item with the key "baz".
	// c.Remove("baz")
}
