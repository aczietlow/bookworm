package bookcache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]CacheObject
	ttl     time.Duration
	mux     *sync.Mutex
}

type CacheObject struct {
	createdAt time.Time
	val       []byte
}

func NewCacheStorage(interval time.Duration) Cache {
	cache := Cache{
		mux:     &sync.Mutex{},
		Entries: make(map[string]CacheObject),
	}
	go cache.reaping(interval)

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, exists := c.Entries[key]; exists == false {
		c.Entries[key] = CacheObject{
			createdAt: time.Now(),
			val:       val,
		}
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if entry, exists := c.Entries[key]; exists == true {
		return entry.val, exists
	}

	return nil, false
}

func (c *Cache) reaping(interval time.Duration) {
	ticker := time.NewTimer(interval)
	defer ticker.Stop()

	// For every ticker.c duration, send a single to reap cache entries
	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for key, entry := range c.Entries {
		if time.Now().Sub(entry.createdAt) > c.ttl {
			delete(c.Entries, key)
		}
	}
}
