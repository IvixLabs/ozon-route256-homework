package inmemory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/cache"
	"testing"
	"time"
)

type testCacheItem struct {
	A string
	B int
}

func TestNewObjectCache(t *testing.T) {
	ctx := context.Background()

	key := "abc"
	item := testCacheItem{A: "aaa", B: 999}
	emptyItem := testCacheItem{}

	cacheObj := NewCache[testCacheItem](emptyItem)

	notFoundItem, err := cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)

	err = cacheObj.Set(ctx, key, item, time.Second)
	assert.NoError(t, err)

	foundItem, err := cacheObj.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, item, foundItem)

	err = cacheObj.Delete(ctx, key)
	assert.NoError(t, err)

	notFoundItem, err = cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)
}

func TestNewObjectCache_Expiration(t *testing.T) {
	ctx := context.Background()

	key := "abc"
	item := testCacheItem{A: "aaa", B: 999}
	emptyItem := testCacheItem{}

	cacheObj := NewCache[testCacheItem](emptyItem)

	err := cacheObj.Set(ctx, key, item, time.Second/5)
	assert.NoError(t, err)

	foundItem, err := cacheObj.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, item, foundItem)

	time.Sleep(time.Second / 2)

	notFoundItem, err := cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)
}

func TestNewObjectCache_NextCache(t *testing.T) {
	ctx := context.Background()

	key := "abc"
	item := testCacheItem{A: "aaa", B: 999}
	emptyItem := testCacheItem{}

	cacheObj := NewCache[testCacheItem](emptyItem)

	emptyNextItem := NextCacheItem[testCacheItem]{}
	nextCache := NewCache[NextCacheItem[testCacheItem]](emptyNextItem)
	cacheObj.SetNext(nextCache)

	notFoundItem, err := cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)

	notFoundNextItem, err := nextCache.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyNextItem, notFoundNextItem)

	err = cacheObj.Set(ctx, key, item, time.Second)
	assert.NoError(t, err)

	foundItem, err := cacheObj.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, item, foundItem)

	foundNextItem, err := nextCache.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, item, foundNextItem.Value)

	err = cacheObj.Delete(ctx, key)
	assert.NoError(t, err)

	notFoundItem, err = cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)

	notFoundNextItem, err = nextCache.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyNextItem, notFoundNextItem)
}

func TestNewObjectCache_NextCacheExpiration(t *testing.T) {
	ctx := context.Background()

	key := "abc"
	item := testCacheItem{A: "aaa", B: 999}
	emptyItem := testCacheItem{}

	cacheObj := NewCache[testCacheItem](emptyItem)

	emptyNextItem := NextCacheItem[testCacheItem]{}
	nextCache := NewCache[NextCacheItem[testCacheItem]](emptyNextItem)
	cacheObj.SetNext(nextCache)

	err := cacheObj.Set(ctx, key, item, time.Second/4)
	assert.NoError(t, err)

	time.Sleep(time.Second / 3)
	foundItem, err := cacheObj.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, item, foundItem)

	time.Sleep(time.Second / 3)

	notFoundItem, err := cacheObj.Get(ctx, key)
	assert.ErrorIs(t, err, cache.ErrItemNotFound)
	assert.Equal(t, emptyItem, notFoundItem)

}
