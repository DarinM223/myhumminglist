package main

import (
	"errors"
	"fmt"
)

// keyValNode is a doubly linked list node
type keyValNode struct {
	key  int
	val  interface{}
	next *keyValNode
	prev *keyValNode
}

// newKeyValNode creates a new key-value node
func newKeyValNode(key int, val interface{}) *keyValNode {
	return &keyValNode{
		key:  key,
		val:  val,
		next: nil,
		prev: nil,
	}
}

// LRUMap is a map that allows you to iterate
// over the entire map in least recently used order
type LRUMap struct {
	dict  map[int]*keyValNode
	front *keyValNode
	rear  *keyValNode
}

// NewLRUMap creates a new LRU map
func NewLRUMap() *LRUMap {
	return &LRUMap{
		dict:  make(map[int]*keyValNode),
		front: nil,
		rear:  nil,
	}
}

// removeFromQueue removes a node from the least recently used linked list
func (lru *LRUMap) removeFromQueue(node *keyValNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		lru.rear = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		lru.front = node.prev
	}
}

// addToFrontOfQueue puts the node as the most recently used
func (lru *LRUMap) addToFrontOfQueue(node *keyValNode) {
	node.next = nil
	node.prev = lru.front

	if lru.rear == nil {
		// node is only node in the queue
		lru.rear = node
	} else {
		// link "front" to the new front
		lru.front.next = node
	}

	// move front pointer to the new front node
	lru.front = node
}

// Add adds a new key-value mapping to the map
func (lru *LRUMap) Add(key int, data interface{}) {
	newNode := newKeyValNode(key, data)
	if node, ok := lru.dict[key]; ok {
		lru.removeFromQueue(node)
		lru.dict[key] = newNode
		lru.addToFrontOfQueue(newNode)
	} else {
		lru.addToFrontOfQueue(newNode)
		lru.dict[key] = newNode
	}
}

// Get retrieves a value from a integer key
// this does not set the node as most recently used, it
// just retrieves the value like a regular map
func (lru *LRUMap) Get(key int) (interface{}, error) {
	if node, ok := lru.dict[key]; ok {
		return node.val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Key %d is not in the LRUMap", key))
	}
}

// Remove removes a key from the map
func (lru *LRUMap) Remove(key int) {
	if node, ok := lru.dict[key]; ok {
		lru.removeFromQueue(node)
		delete(lru.dict, key)
	}
}

// Keys returns an array of indexes from least recently used
// to most recently used
func (lru *LRUMap) Keys() []int {
	var keys []int
	for node := lru.rear; node != nil; node = node.next {
		keys = append(keys, node.key)
	}
	return keys
}
