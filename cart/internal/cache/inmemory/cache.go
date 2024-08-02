package inmemory

import (
	"context"
	"errors"
	"route256/cart/internal/cache"
	"sync"
	"time"
)

type NextCacheItem[V any] struct {
	Value      V
	Expiration time.Duration
}

type Cache[V any] struct {
	cache     map[string]V
	cacheMx   sync.RWMutex
	emptyItem V
	nextCache cache.Cache[NextCacheItem[V]]
}

func NewCache[V any](emptyItem V) *Cache[V] {
	return &Cache[V]{cache: make(map[string]V), emptyItem: emptyItem}
}

func (c *Cache[V]) Set(ctx context.Context, k string, v V, expiration time.Duration) error {

	c.set(ctx, k, v, expiration)

	if c.nextCache != nil {
		err := c.nextCache.Set(ctx, k, NextCacheItem[V]{Value: v, Expiration: expiration}, expiration*2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache[V]) set(_ context.Context, k string, v V, expiration time.Duration) {
	c.cacheMx.Lock()
	c.cache[k] = v
	c.cacheMx.Unlock()

	time.AfterFunc(expiration, func() {
		c.cacheMx.Lock()
		delete(c.cache, k)
		c.cacheMx.Unlock()
	})
}

func (c *Cache[V]) Get(ctx context.Context, k string) (V, error) {

	c.cacheMx.RLock()
	item, ok := c.cache[k]
	c.cacheMx.RUnlock()

	if !ok {

		if c.nextCache != nil {
			ncItem, err := c.nextCache.Get(ctx, k)
			if err != nil {
				if !errors.Is(err, cache.ErrItemNotFound) {
					return c.emptyItem, err
				}
			} else {
				c.set(ctx, k, ncItem.Value, ncItem.Expiration)
				return ncItem.Value, nil

			}
		}

		return c.emptyItem, cache.ErrItemNotFound
	}

	return item, nil
}

func (c *Cache[V]) Delete(ctx context.Context, k string) error {
	c.cacheMx.Lock()
	delete(c.cache, k)
	c.cacheMx.Unlock()

	if c.nextCache != nil {
		err := c.nextCache.Delete(ctx, k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache[V]) SetNext(next cache.Cache[NextCacheItem[V]]) {
	c.nextCache = next
}
