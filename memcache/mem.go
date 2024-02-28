package memcache

import (
	"context"
	"sync"
	"time"
)

type caching struct {
	store  map[string]interface{}
	locker *sync.RWMutex
}

func NewCaching() *caching {
	return &caching{
		store:  make(map[string]interface{}),
		locker: new(sync.RWMutex),
	}
}

func (c *caching) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.store[key] = value

	return nil
}

func (c *caching) Get(ctx context.Context, key string, value interface{}) error {
	c.locker.RLock()
	defer c.locker.RUnlock()
	value = c.store[key]

	return nil
}

func (c *caching) Delete(ctx context.Context, key string) error {
	c.locker.Lock()
	defer c.locker.Unlock()
	delete(c.store, key)
	return nil
}
