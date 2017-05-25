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
	cache := New(3)
	cache.addNode(NewNode("hoge", 111))
	cache.addNode(NewNode("fuga", 222))
	cache.addNode(NewNode("fooo", 333))

	if size := cache.CurrentSize(); size != 3 {
		t.Error("Expected 3 but got", size)
	}

	keys1 := cache.AllKeys()
	keys2 := cache.AllKeysReversed()

	if keys1[0] != keys2[2] {
		t.Error("Expected equal but got", keys1[0], keys2[2])
	}
	if keys1[1] != keys2[1] {
		t.Error("Expected equal but got", keys1[1], keys2[1])
	}
	if keys1[2] != keys2[0] {
		t.Error("Expected equal but got", keys1[2], keys2[0])
	}
}

func TestRemoveNode(t *testing.T) {
	n1 := NewNode("hoge", 111)
	n2 := NewNode("fuga", 222)
	n3 := NewNode("fooo", 333)
	n4 := NewNode("baar", 444)

	cache := New(4)
	cache.addNode(n1)
	cache.addNode(n2)
	cache.addNode(n3)
	cache.addNode(n4)

	if size := cache.CurrentSize(); size != 4 {
		t.Error("Expected 4 but got", size)
	}

	cache.removeNode(n2)

	if size := cache.CurrentSize(); size != 3 {
		t.Error("Expected 3 but got", size)
	}

	keys1 := cache.AllKeys()
	keys2 := cache.AllKeysReversed()

	if keys1[0] != keys2[2] {
		t.Error("Expected equal but got", keys1[0], keys2[2])
	}
	if keys1[1] != keys2[1] {
		t.Error("Expected equal but got", keys1[1], keys2[1])
	}
	if keys1[2] != keys2[0] {
		t.Error("Expected equal but got", keys1[2], keys2[0])
	}

	cache.removeNode(n1)

	if size := cache.CurrentSize(); size != 2 {
		t.Error("Expected 2 but got", size)
	}

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	if keys1[0] != keys2[1] {
		t.Error("Expected equal but got", keys1[0], keys2[1])
	}
	if keys1[1] != keys2[0] {
		t.Error("Expected equal but got", keys1[1], keys2[0])
	}

	cache.removeNode(n4)

	if size := cache.CurrentSize(); size != 1 {
		t.Error("Expected 1 but got", size)
	}

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	if keys1[0] != keys2[0] {
		t.Error("Expected equal but got", keys1[0], keys2[0])
	}

	cache.removeNode(n3)

	if size := cache.CurrentSize(); size != 0 {
		t.Error("Expected 0 but got", size)
	}

	keys1 = cache.AllKeys()
	keys2 = cache.AllKeysReversed()

	if len(keys1) != 0 {
		t.Error("Expected no keys but got", keys1)
	}
	if len(keys2) != 0 {
		t.Error("Expected no keys but got", keys2)
	}
}

func TestSet(t *testing.T) {
	cache := New(3)
	cache.Set("hoge", 111)
	cache.Set("fuga", 222)
	cache.Set("fooo", 333)
	cache.Set("baar", 444)

	if size := cache.CurrentSize(); size != 3 {
		t.Error("Expected 3 but got", size)
	}

	var v interface{}
	var e error

	v, e = cache.Get("hoge")

	if e == nil {
		t.Error("Expected error but got", e)
	}
	if v != nil {
		t.Error("Expected nil but got", v)
	}

	v, e = cache.Get("fuga")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if v != 222 {
		t.Error("Expected 222 but got", v)
	}

	v, e = cache.Get("fooo")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if v != 333 {
		t.Error("Expected 333 but got", v)
	}

	v, e = cache.Get("baar")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if v != 444 {
		t.Error("Expected 444 but got", v)
	}
}

func TestGet(t *testing.T) {
	cache := New(3)
	cache.Set("hoge", 111)
	cache.Set("fuga", 222)
	cache.Set("fooo", 333)

	cache.Get("hoge") // "fuga" is the head now
	cache.Get("fuga") // "fooo" is the head now

	cache.Set("baar", 444) // "fooo" is deleted

	if size := cache.CurrentSize(); size != 3 {
		t.Error("Expected 3 but got", size)
	}

	var v interface{}
	var e error

	v, e = cache.Get("fooo")

	if v != nil {
		t.Error("Expected nil but got", v)
	}
	if e == nil {
		t.Error("Expected error but got", e)
	}

	v, e = cache.Get("hoge")

	if v != 111 {
		t.Error("Expected 111 but got", v)
	}
	if e != nil {
		t.Error("Expected nil but got", e)
	}

	v, e = cache.Get("fuga")

	if v != 222 {
		t.Error("Expected 222 but got", v)
	}
	if e != nil {
		t.Error("Expected nil but got", e)
	}

	v, e = cache.Get("baar")

	if v != 444 {
		t.Error("Expected 444 but got", v)
	}
	if e != nil {
		t.Error("Expected nil but got", e)
	}
}

func TestDelete(t *testing.T) {
	cache := New(3)
	cache.Set("hoge", 111)
	cache.Set("fuga", 222)
	cache.Set("fooo", 333)

	var e error

	e = cache.Delete("fuga")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if s := cache.CurrentSize(); s != 2 {
		t.Error("Expected 2 but got", s)
	}

	e = cache.Delete("fooo")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if s := cache.CurrentSize(); s != 1 {
		t.Error("Expected 1 but got", s)
	}

	e = cache.Delete("hoge")

	if e != nil {
		t.Error("Expected nil but got", e)
	}
	if s := cache.CurrentSize(); s != 0 {
		t.Error("Expected 0 but got", s)
	}

	e = cache.Delete("nonexisting")

	if e == nil {
		t.Error("Expected error but got", e)
	}
}
