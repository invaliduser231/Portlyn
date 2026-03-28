package proxy

import (
	"container/list"
	"sync"
	"time"
)

type ttlLRU[K comparable, V any] struct {
	mu         sync.Mutex
	ll         *list.List
	cache      map[K]*list.Element
	capacity   int
	defaultTTL time.Duration
}

type ttlEntry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

func newTTLLRU[K comparable, V any](capacity int, ttl time.Duration) *ttlLRU[K, V] {
	if capacity <= 0 {
		capacity = 1
	}
	return &ttlLRU[K, V]{
		ll:         list.New(),
		cache:      make(map[K]*list.Element, capacity),
		capacity:   capacity,
		defaultTTL: ttl,
	}
}

func (c *ttlLRU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	element, ok := c.cache[key]
	if !ok {
		return zero, false
	}
	entry := element.Value.(*ttlEntry[K, V])
	if !entry.expiresAt.IsZero() && time.Now().UTC().After(entry.expiresAt) {
		c.removeElement(element)
		return zero, false
	}
	c.ll.MoveToFront(element)
	return entry.value, true
}

func (c *ttlLRU[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.cache[key]; ok {
		entry := element.Value.(*ttlEntry[K, V])
		entry.value = value
		entry.expiresAt = time.Now().UTC().Add(c.defaultTTL)
		c.ll.MoveToFront(element)
		return
	}

	entry := &ttlEntry[K, V]{
		key:       key,
		value:     value,
		expiresAt: time.Now().UTC().Add(c.defaultTTL),
	}
	element := c.ll.PushFront(entry)
	c.cache[key] = element

	if c.ll.Len() > c.capacity {
		c.removeElement(c.ll.Back())
	}
}

func (c *ttlLRU[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.cache[key]; ok {
		c.removeElement(element)
	}
}

func (c *ttlLRU[K, V]) Purge() {
	c.mu.Lock()
	c.ll.Init()
	c.cache = make(map[K]*list.Element, c.capacity)
	c.mu.Unlock()
}

func (c *ttlLRU[K, V]) removeElement(element *list.Element) {
	if element == nil {
		return
	}
	c.ll.Remove(element)
	entry := element.Value.(*ttlEntry[K, V])
	delete(c.cache, entry.key)
}
