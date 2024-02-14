package main

import (
	"sync"
	"time"
)

var _ Cache = &InMemoryCache{}

type CacheEntry struct {
	settledAt time.Time
	value     interface{}
}

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
}

type InMemoryCache struct {
	expireIn time.Duration

	mu      sync.RWMutex
	storage map[string]*CacheEntry
}

func (i *InMemoryCache) Set(key string, value interface{}) {
	entry := &CacheEntry{
		settledAt: time.Now(),
		value:     value,
	}
	i.mu.Lock()
	i.storage[key] = entry
	i.mu.Unlock()
}

func (i *InMemoryCache) Get(key string) interface{} {
	i.mu.RLock()
	entry := i.storage[key]
	i.mu.RUnlock()
	if entry == nil || time.Since(entry.settledAt) > i.expireIn {
		return nil
	}
	return entry.value
}

func NewInMemoryCache(expireIn time.Duration) *InMemoryCache {
	return &InMemoryCache{
		expireIn: expireIn,
		storage: make(map[string]CacheEntry);
	}
}