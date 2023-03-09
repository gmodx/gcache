package memcache

import (
	"testing"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ColumnA string
}

func Test_Cache(t *testing.T) {
	c := NewDefault[testItem]()

	key1 := "key1"
	val1 := testItem{
		ColumnA: "val1",
	}

	key2 := "key2"
	val2 := testItem{
		ColumnA: "val2",
	}

	resultVal1 := c.Get(key1)
	assert.Nil(t, resultVal1)

	c.Set(key1, val1, abstract.CacheEntryOptions{Expiration: time.Hour})
	c.Set(key2, val2, abstract.CacheEntryOptions{Expiration: time.Hour})

	resultVal1 = c.Get(key1)
	assert.EqualValues(t, val1, *resultVal1)
	resultVal2 := c.Get(key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Remove(key1)
	resultVal1 = c.Get(key1)
	assert.Nil(t, resultVal1)
	resultVal2 = c.Get(key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Set(key1, val1, abstract.CacheEntryOptions{Expiration: time.Hour})
	resultVal1 = c.Get(key1)
	assert.EqualValues(t, val1, *resultVal1)
	resultVal2 = c.Get(key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Flush()
	resultVal1 = c.Get(key1)
	assert.Nil(t, resultVal1)
	resultVal2 = c.Get(key2)
	assert.Nil(t, resultVal2)
}

func Test_CacheTime(t *testing.T) {
	tc := New[int](MemCacheOptions{CleanupInterval: 3 * time.Millisecond})
	tc.Set("b", 2, abstract.CacheEntryOptions{Expiration: abstract.NoExpiration})
	tc.Set("c", 3, abstract.CacheEntryOptions{Expiration: 20 * time.Millisecond})
	tc.Set("d", 4, abstract.CacheEntryOptions{Expiration: 70 * time.Millisecond})

	<-time.After(25 * time.Millisecond)
	val := tc.Get("c")
	if val != nil {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(30 * time.Millisecond)
	val = tc.Get("b")
	if val == nil {
		t.Error("Did not find b even though it was set to never expire")
	}

	val = tc.Get("d") // 25+30==55, 55<70
	if val == nil {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(20 * time.Millisecond) // 25+30+20==75, 75>70
	val = tc.Get("d")
	if val != nil {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}
