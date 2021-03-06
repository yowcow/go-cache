package lrucache

import (
	"fmt"
	"github.com/yowcow/go-cache/cache"
	"sync"
)

type LRUCacheNode struct {
	key  string
	val  interface{}
	prev *LRUCacheNode
	next *LRUCacheNode
}

type LRUCache struct {
	maxSize     int64
	currentSize int64
	head        *LRUCacheNode
	tail        *LRUCacheNode
	keyMap      map[string]*LRUCacheNode
	mutex       *sync.Mutex
}

type LRUCacher interface {
	cache.Cacher
	addNode(*LRUCacheNode)
	removeNode(*LRUCacheNode)
}

func New(maxSize int64) LRUCacher {
	return &LRUCache{
		maxSize:     maxSize,
		currentSize: 0,
		head:        nil,
		tail:        nil,
		keyMap:      map[string]*LRUCacheNode{},
		mutex:       &sync.Mutex{},
	}
}

func NewNode(key string, val interface{}) *LRUCacheNode {
	return &LRUCacheNode{
		key:  key,
		val:  val,
		prev: nil,
		next: nil,
	}
}

func (self *LRUCache) MaxSize() int64 {
	return self.maxSize
}

func (self *LRUCache) CurrentSize() int64 {
	return self.currentSize
}

func (self *LRUCache) AllKeys() []string {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	node := self.head
	keys := make([]string, self.CurrentSize())
	for i := 0; node != nil; i++ {
		keys[i] = node.key
		node = node.next
	}
	return keys
}

func (self *LRUCache) AllKeysReversed() []string {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	node := self.tail
	keys := make([]string, self.CurrentSize())
	for i := 0; node != nil; i++ {
		keys[i] = node.key
		node = node.prev
	}
	return keys
}

func (self *LRUCache) Set(key string, val interface{}) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if node := self.keyMap[key]; node != nil { // Existing key
		node.val = val
		self.removeNode(node)
		self.addNode(node)
	} else { // New key
		if self.CurrentSize() == self.MaxSize() {
			head := self.head
			delete(self.keyMap, head.key)
			self.removeNode(head)
		}
		node = NewNode(key, val)
		self.addNode(node)
		self.keyMap[key] = node
	}
	return nil
}

func (self *LRUCache) Get(key string) (interface{}, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if node := self.keyMap[key]; node != nil {
		if node != self.tail { // The node is not the tail
			self.removeNode(node)
			self.addNode(node)
		}
		return node.val, nil
	}
	return nil, fmt.Errorf("Key %s does not exist", key)
}

func (self *LRUCache) Delete(key string) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if node := self.keyMap[key]; node != nil {
		self.removeNode(node)
		delete(self.keyMap, key)
		return nil
	}
	return fmt.Errorf("Key %s does not exist", key)
}

func (self *LRUCache) addNode(node *LRUCacheNode) {
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

func (self *LRUCache) removeNode(node *LRUCacheNode) {
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
