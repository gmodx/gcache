package memcache

import (
	"context"
	"testing"
	"time"

	"github.com/gmodx/gcache/abstract"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ColumnA string
}

func Test_Cache(t *testing.T) {
	ctx := context.TODO()

	c := NewDefault[testItem]()

	key1 := "key1"
	val1 := testItem{
		ColumnA: "val1",
	}

	key2 := "key2"
	val2 := testItem{
		ColumnA: "val2",
	}

	resultVal1, _ := c.Get(ctx, key1)
	assert.Nil(t, resultVal1)

	c.Set(ctx, key1, val1, abstract.CacheEntryOptions{Expiration: time.Hour})
	c.Set(ctx, key2, val2, abstract.CacheEntryOptions{Expiration: time.Hour})

	resultVal1, _ = c.Get(ctx, key1)
	assert.EqualValues(t, val1, *resultVal1)
	resultVal2, _ := c.Get(ctx, key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Remove(ctx, key1)
	resultVal1, _ = c.Get(ctx, key1)
	assert.Nil(t, resultVal1)
	resultVal2, _ = c.Get(ctx, key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Set(ctx, key1, val1, abstract.CacheEntryOptions{Expiration: time.Hour})
	resultVal1, _ = c.Get(ctx, key1)
	assert.EqualValues(t, val1, *resultVal1)
	resultVal2, _ = c.Get(ctx, key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Flush()
	resultVal1, _ = c.Get(ctx, key1)
	assert.Nil(t, resultVal1)
	resultVal2, _ = c.Get(ctx, key2)
	assert.Nil(t, resultVal2)
}

func Test_CacheTime(t *testing.T) {
	ctx := context.TODO()
	tc := New[int](Options{CleanupInterval: 3 * time.Millisecond})
	tc.Set(ctx, "b", 2, abstract.CacheEntryOptions{Expiration: abstract.NoExpiration})
	tc.Set(ctx, "c", 3, abstract.CacheEntryOptions{Expiration: 20 * time.Millisecond})
	tc.Set(ctx, "d", 4, abstract.CacheEntryOptions{Expiration: 70 * time.Millisecond})

	<-time.After(25 * time.Millisecond)
	val, _ := tc.Get(ctx, "c")
	if val != nil {
		t.Error("Found c when it should have been automatically deleted")
	}

	<-time.After(30 * time.Millisecond)
	val, _ = tc.Get(ctx, "b")
	if val == nil {
		t.Error("Did not find b even though it was set to never expire")
	}

	val, _ = tc.Get(ctx, "d") // 25+30==55, 55<70
	if val == nil {
		t.Error("Did not find d even though it was set to expire later than the default")
	}

	<-time.After(20 * time.Millisecond) // 25+30+20==75, 75>70
	val, _ = tc.Get(ctx, "d")
	if val != nil {
		t.Error("Found d when it should have been automatically deleted (later than the default)")
	}
}
