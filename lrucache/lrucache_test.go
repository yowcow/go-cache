package lrucache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	cache := New(5)

	assert.Equal(int64(5), cache.MaxSize())
	assert.Equal(int64(0), cache.CurrentSize())
	assert.True(assert.EqualValues([]string{}, cache.AllKeys()))
}

func TestAddNode(t *testing.T) {
	assert := assert.New(t)

	cache := New(3)
	cache.addNode(NewNode("hoge", 111))
	cache.addNode(NewNode("fuga", 222))
	cache.addNode(NewNode("fooo", 333))

	assert.Equal(int64(3), cache.CurrentSize())

	keys1 := cache.AllKeys()
	keys2 := cache.AllKeysReversed()

	assert.Equal(keys1[0], keys2[2])
	assert.Equal(keys1[1], keys2[1])
	assert.Equal(keys1[2], keys2[0])
}

func TestRemoveNode(t *testing.T) {
	assert := assert.New(t)

	n1 := NewNode("hoge", 111)
	n2 := NewNode("fuga", 222)
	n3 := NewNode("fooo", 333)
	n4 := NewNode("baar", 444)

	cache := New(4)
	cache.addNode(n1)
	cache.addNode(n2)
	cache.addNode(n3)
	cache.addNode(n4)

	assert.Equal(int64(4), cache.CurrentSize())

	cache.removeNode(n2)

	assert.Equal(int64(3), cache.CurrentSize())

	keys1 := cache.AllKeys()
	keys2 := cache.AllKeysReversed()

	assert.Equal(keys1[0], keys2[2])
	assert.Equal(keys1[1], keys2[1])
	assert.Equal(keys1[2], keys2[0])

	cache.removeNode(n1)

	assert.Equal(int64(2), cache.CurrentSize())

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	assert.Equal(keys1[0], keys2[1])
	assert.Equal(keys1[1], keys2[0])

	cache.removeNode(n4)

	assert.Equal(int64(1), cache.CurrentSize())

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	assert.Equal(keys1[0], keys2[0])

	cache.removeNode(n3)

	assert.Equal(int64(0), cache.CurrentSize())

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	assert.Equal(0, len(keys1))
	assert.Equal(0, len(keys2))
}

func TestSet(t *testing.T) {
	var v interface{}
	var e error

	assert := assert.New(t)

	cache := New(3)
	cache.Set("hoge", 11)
	cache.Set("hoge", 111)

	assert.Equal(int64(1), cache.CurrentSize())

	v, e = cache.Get("hoge")

	assert.Equal(nil, e)
	assert.Equal(111, v)

	cache.Set("fuga", 222)
	cache.Set("fooo", 333)
	cache.Set("baar", 444)

	assert.Equal(int64(3), cache.CurrentSize())

	v, e = cache.Get("hoge")

	assert.NotEqual(nil, e)
	assert.Equal(nil, v)

	v, e = cache.Get("fuga")

	assert.Equal(nil, e)
	assert.Equal(222, v)

	v, e = cache.Get("fooo")

	assert.Equal(nil, e)
	assert.Equal(333, v)

	v, e = cache.Get("baar")

	assert.Equal(nil, e)
	assert.Equal(444, v)
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	cache := New(3)
	cache.Set("hoge", 111)
	cache.Set("fuga", 222)
	cache.Set("fooo", 333)

	cache.Get("hoge") // "fuga" is the head now
	cache.Get("fuga") // "fooo" is the head now

	cache.Set("baar", 444) // "fooo" is deleted

	assert.Equal(int64(3), cache.CurrentSize())

	var v interface{}
	var e error

	v, e = cache.Get("fooo")

	assert.Equal(nil, v)
	assert.NotEqual(nil, e)

	v, e = cache.Get("hoge")

	assert.Equal(111, v)
	assert.Equal(nil, e)

	v, e = cache.Get("fuga")

	assert.Equal(222, v)
	assert.Equal(nil, e)

	v, e = cache.Get("baar")

	assert.Equal(444, v)
	assert.Equal(nil, e)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	cache := New(3)
	cache.Set("hoge", 111)
	cache.Set("fuga", 222)
	cache.Set("fooo", 333)

	var e error

	e = cache.Delete("fuga")

	assert.Equal(nil, e)
	assert.Equal(int64(2), cache.CurrentSize())

	e = cache.Delete("fooo")

	assert.Equal(nil, e)
	assert.Equal(int64(1), cache.CurrentSize())

	e = cache.Delete("hoge")

	assert.Equal(nil, e)
	assert.Equal(int64(0), cache.CurrentSize())

	e = cache.Delete("nonexisting")

	assert.NotEqual(nil, e)
}
