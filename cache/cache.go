package cache

import (
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *Lru
	maxBytes   int64
}

func (c *cache) get(key string) (ByteView, bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = New(c.maxBytes, nil)
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}

	return ByteView{}, false
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = New(c.maxBytes, nil)
	}

	c.lru.Add(key, value)
}