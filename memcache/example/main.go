package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/gmodx/gcache/memcache"
)

func main() {
	var c abstract.ICache[string] = memcache.New[string](memcache.Options{CleanupInterval: 10 * time.Second})

	ctx := context.TODO()
	// Set a value in the cache with the key "baz" and the value "bar", specifying that it has no expiration.
	c.Set(ctx, "baz", "bar", abstract.CacheEntryOptions{Expiration: abstract.NoExpiration})

	// Get the value stored in the cache with the key "baz".
	val, _ := c.Get(ctx, "baz")
	if val != nil {
		fmt.Println(val)
	}

	// Refresh the expiration time of the cached item with the key "baz".
	c.Refresh(ctx, "baz")

	// Remove the cached item with the key "baz".
	c.Remove(ctx, "baz")
}
