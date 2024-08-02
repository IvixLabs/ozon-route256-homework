package noop

import (
	"context"
	"log"
	"route256/cart/internal/cache"
	"time"
)

type Cache[V any] struct {
	emptyItem V
}

func NewCache[V any](emptyItem V) *Cache[V] {
	return &Cache[V]{emptyItem: emptyItem}
}

func (c *Cache[V]) Set(_ context.Context, _ string, _ V, _ time.Duration) error {

	return nil
}

func (c *Cache[V]) Get(_ context.Context, _ string) (V, error) {
	log.Println("NOOP CACHE GET")

	return c.emptyItem, cache.ErrItemNotFound
}

func (c *Cache[V]) Delete(_ context.Context, _ string) error {
	return nil
}
