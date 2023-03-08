package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	ColumnA string
}

func Test_memcache(t *testing.T) {
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

	c.Set(key1, val1)
	c.Set(key2, val2)

	resultVal1 = c.Get(key1)
	assert.EqualValues(t, val1, *resultVal1)
	resultVal2 := c.Get(key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Remove(key1)
	resultVal1 = c.Get(key1)
	assert.Nil(t, resultVal1)
	resultVal2 = c.Get(key2)
	assert.EqualValues(t, val2, *resultVal2)

	c.Set(key1, val1)
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
