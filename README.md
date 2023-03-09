gcache
========
[![GoDoc](https://godoc.org/github.com/gmodx/gcache?status.svg)](https://pkg.go.dev/github.com/gmodx/gcache)

An golang cache library, support generic.

## Installation

`go get github.com/gmodx/gcache`

## Usage

### memory cache example
```go
package main

import (
	"fmt"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/gmodx/gcache/memcache"
)

func main() {
	var c abstract.ICache[string] = memcache.New[string](memcache.MemCacheOptions{CleanupInterval: 10 * time.Second})

	// Set a value in the cache with the key "baz" and the value "bar", specifying that it has no expiration.
	c.Set("baz", "bar", abstract.CacheEntryOptions{Expiration: abstract.NoExpiration})

	// Get the value stored in the cache with the key "baz".
	val := c.Get("baz")
	if val != nil {
		fmt.Println(val)
	}

	// Refresh the expiration time of the cached item with the key "baz".
	c.Refresh("baz")

	// Remove the cached item with the key "baz".
	c.Remove("baz")
}
```