package cache

import (
	"runtime"
	"sync"
	"time"
)

const (
	defaultCleanupInterval time.Duration = 10 * time.Minute
)

type Cacher interface {
	Set(key string, value any, ttl time.Duration)
	Get(key string) (any, bool)
	Delete(key string) bool
	DeleteExpired()
	Flush()
}

type Cache struct {
	*cache
}

type cache struct {
	sync.Mutex
	ttl     time.Duration
	items   map[string]*item
	clearer *clearJob
}

var _ Cacher = (*Cache)(nil)

func New(ttl time.Duration) (*Cache, error) {
	c := newCache(ttl)

	C := &Cache{c}

	runClearer(c, defaultCleanupInterval)
	runtime.SetFinalizer(C, stopClearer)

	return C, nil
}

func newCache(ttl time.Duration) *cache {
	if ttl == 0 {
		ttl = -1
	}
	c := &cache{
		ttl:   ttl,
		items: map[string]*item{},
	}
	return c
}

// Set key to cache -.
func (c *cache) Set(key string, value any, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()

	if ttl == 0 {
		ttl = c.ttl
	}

	c.set(key, value, ttl)
}

// Internal set with logic -.
func (c *cache) set(key string, value any, ttl time.Duration) {
	t := time.Now().Add(ttl)

	c.items[key] = &item{
		Object:     value,
		Expiration: &t,
	}
}

// Get value by key -.
func (c *cache) Get(key string) (any, bool) {
	c.Lock()
	defer c.Unlock()

	x, found := c.get(key)

	return x, found
}

func (c *cache) get(k string) (any, bool) {
	item, found := c.items[k]

	if !found {
		return nil, false
	}

	if item.Expired() {
		c.delete(k)
		return nil, false
	}

	return item.Object, true
}

func (c *cache) Delete(key string) bool {
	c.Lock()
	defer c.Unlock()

	_, found := c.get(key)

	c.delete(key)

	return found
}

func (c *cache) delete(key string) {
	delete(c.items, key)
}

func (c *cache) DeleteExpired() {
	c.Lock()
	defer c.Unlock()

	for k, v := range c.items {
		if v.Expired() {
			c.delete(k)
		}
	}
}

func (c *cache) Flush() {
	c.Lock()
	defer c.Unlock()

	c.items = map[string]*item{}
}
