// Copyright 2023 enpsl. All rights reserved.

// LRU Cache Manager

package base

import (
	"sync"
)

type DoubleLink struct {
	Val  interface{}
	Key  string
	Next *DoubleLink
	Prev *DoubleLink
}

type LRUCache struct {
	size       int
	capacity   int
	head, tail *DoubleLink
	hashMap    map[string]*DoubleLink
	mu         sync.RWMutex
}

func initNode() *DoubleLink {
	return new(DoubleLink)
}

func CacheConstructor(capacity int) *LRUCache {
	lru := LRUCache{
		capacity: capacity,
		size:     0,
		head:     initNode(),
		tail:     initNode(),
		hashMap:  make(map[string]*DoubleLink, capacity),
	}
	lru.tail.Prev = lru.head
	lru.head.Next = lru.tail
	return &lru
}

func (this *LRUCache) Get(key string) (interface{}, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if node, ok := this.hashMap[key]; ok {
		this.MoveToHead(node)
		return node.Val, true
	} else {
		return nil, false
	}
}

func (this *LRUCache) Put(key string, value interface{}) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if _, ok := this.hashMap[key]; ok {
		node := this.hashMap[key]
		node.Val = value
		this.MoveToHead(node)
	} else {
		node := initNode()
		node.Key = key
		node.Val = value
		this.AddToHead(node)
		this.size++
		this.hashMap[key] = node
		if this.size > this.capacity {
			removed := this.RemoveTail()
			delete(this.hashMap, removed.Key)
			this.size--
		}
	}
}

func (this *LRUCache) AddToHead(node *DoubleLink) {
	node.Prev = this.head
	node.Next = this.head.Next
	this.head.Next.Prev = node
	this.head.Next = node
}

func (this *LRUCache) RemoveNode(node *DoubleLink) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

func (this *LRUCache) RemoveTail() *DoubleLink {
	remove := this.tail.Prev
	this.RemoveNode(remove)
	return remove
}

func (this *LRUCache) MoveToHead(node *DoubleLink) {
	this.RemoveNode(node)
	this.AddToHead(node)
}

func (this *LRUCache) Flush() {
	this.hashMap = make(map[string]*DoubleLink, this.capacity)
	this.size = 0
	this.head = initNode()
	this.tail = initNode()
	this.tail.Prev = this.head
	this.head.Next = this.tail
}
