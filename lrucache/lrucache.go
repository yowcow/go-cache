package lrucache

import (
	"fmt"
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
}

type LRUCacheInterface interface {
	MaxSize() int64
	CurrentSize() int64
	AllKeys() []string
	AllKeysReversed() []string
	AddNode(*LRUCacheNode)
	RemoveNode(*LRUCacheNode)
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	Delete(string) error
}

func New(maxSize int64) LRUCacheInterface {
	return &LRUCache{
		maxSize:     maxSize,
		currentSize: 0,
		head:        nil,
		tail:        nil,
		keyMap:      map[string]*LRUCacheNode{},
	}
}

func NewNode(key string, val interface{}) *LRUCacheNode {
	return &LRUCacheNode{key, val, nil, nil}
}

func (self *LRUCache) MaxSize() int64 {
	return self.maxSize
}

func (self *LRUCache) CurrentSize() int64 {
	return self.currentSize
}

func (self *LRUCache) AllKeys() []string {
	node := self.head
	result := []string{}
	for node != nil {
		result = append(result, node.key)
		node = node.next
	}
	return result
}

func (self *LRUCache) AllKeysReversed() []string {
	node := self.tail
	result := []string{}
	for node != nil {
		result = append(result, node.key)
		node = node.prev
	}
	return result
}

func (self *LRUCache) AddNode(node *LRUCacheNode) {
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

func (self *LRUCache) RemoveNode(node *LRUCacheNode) {
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

func (self *LRUCache) Set(key string, val interface{}) error {
	if node := self.keyMap[key]; node != nil { // Existing key
		self.RemoveNode(node)
		self.AddNode(node)
	} else { // New key
		if self.CurrentSize() == self.MaxSize() {
			head := self.head
			self.Delete(head.key)
		}
		node = NewNode(key, val)
		self.AddNode(node)
		self.keyMap[key] = node
	}
	return nil
}

func (self *LRUCache) Get(key string) (interface{}, error) {
	if node := self.keyMap[key]; node != nil {
		if node != self.tail { // The node is not the tail
			self.RemoveNode(node)
			self.AddNode(node)
		}
		return node.val, nil
	}
	return nil, fmt.Errorf("Key %s does not exist", key)
}

func (self *LRUCache) Delete(key string) error {
	if node := self.keyMap[key]; node != nil {
		self.RemoveNode(node)
		delete(self.keyMap, key)
		return nil
	}
	return fmt.Errorf("Key %s does not exist", key)
}
