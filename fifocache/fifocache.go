package fifocache

import (
	"fmt"
	"github.com/yowcow/go-cache/cache"
)

type FIFOCacheNode struct {
	key  string
	val  interface{}
	prev *FIFOCacheNode
	next *FIFOCacheNode
}

type FIFOCache struct {
	maxSize     int64
	currentSize int64
	head        *FIFOCacheNode
	tail        *FIFOCacheNode
	keyMap      map[string]*FIFOCacheNode
}

type FIFOCacher interface {
	cache.Cacher
	addNode(*FIFOCacheNode)
	removeNode(*FIFOCacheNode)
}

func New(maxSize int64) FIFOCacher {
	return &FIFOCache{
		maxSize:     maxSize,
		currentSize: 0,
		head:        nil,
		tail:        nil,
		keyMap:      map[string]*FIFOCacheNode{},
	}
}

func NewNode(key string, val interface{}) *FIFOCacheNode {
	return &FIFOCacheNode{
		key:  key,
		val:  val,
		prev: nil,
		next: nil,
	}
}

func (self *FIFOCache) MaxSize() int64 {
	return self.maxSize
}

func (self *FIFOCache) CurrentSize() int64 {
	return self.currentSize
}

func (self *FIFOCache) AllKeys() []string {
	node := self.head
	keys := make([]string, self.CurrentSize())
	for i := 0; node != nil; i++ {
		keys[i] = node.key
		node = node.next
	}
	return keys
}

func (self *FIFOCache) AllKeysReversed() []string {
	node := self.tail
	keys := make([]string, self.CurrentSize())
	for i := 0; node != nil; i++ {
		keys[i] = node.key
		node = node.prev
	}
	return keys
}

func (self *FIFOCache) Set(key string, val interface{}) error {
	if node := self.keyMap[key]; node != nil { // Existing key
		node.val = val
	} else { // New key
		if self.CurrentSize() == self.MaxSize() {
			head := self.head
			self.Delete(head.key)
		}
		node = NewNode(key, val)
		self.addNode(node)
		self.keyMap[key] = node
	}
	return nil
}

func (self *FIFOCache) Get(key string) (interface{}, error) {
	if node := self.keyMap[key]; node != nil {
		return node.val, nil
	}
	return nil, fmt.Errorf("Key %s does not exist", key)
}

func (self *FIFOCache) Delete(key string) error {
	if node := self.keyMap[key]; node != nil {
		self.removeNode(node)
		delete(self.keyMap, key)
		return nil
	}
	return fmt.Errorf("Key %s does not exist", key)
}

func (self *FIFOCache) addNode(node *FIFOCacheNode) {
	if tail := self.tail; tail != nil {
		tail.next = node
		node.prev = tail
		self.tail = node
	} else {
		self.head = node
		self.tail = node
	}

	self.currentSize += 1
}

func (self *FIFOCache) removeNode(node *FIFOCacheNode) {
	if node == self.head && node == self.tail { // Removing the last node
		self.head = nil
		self.tail = nil
	} else if node == self.head { // Removing the head node
		next := node.next
		next.prev = nil
		self.head = next
	} else if node == self.tail { // Removing the tail node
		prev := node.prev
		prev.next = nil
		self.tail = prev
	} else { // Removing one in the middle
		prev := node.prev
		next := node.next
		prev.next = next
		next.prev = prev
	}

	self.currentSize -= 1
}
