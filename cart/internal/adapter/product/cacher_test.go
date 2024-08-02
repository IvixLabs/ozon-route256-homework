package product

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"route256/cart/internal/cache/inmemory"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
	"route256/cart/internal/usecase/cart/mock"
	"sync"
	"testing"
	"time"
)

func TestCacher(t *testing.T) {
	ctx := context.Background()
	provider := mock.NewProductProviderMock(t)

	sku := model.Sku(123)
	product := &model.Product{Name: "Prod123", Price: 1230}

	skuNotExists := model.Sku(111)

	provider.
		GetMock.When(ctx, sku).Then(product, nil).
		GetMock.When(ctx, skuNotExists).Then(nil, cart.ErrProductSkuNotFound)

	inmemoryCache := inmemory.NewCache[CacheItem](CacheItem{})
	cacher := NewCacher(provider, inmemoryCache)

	retProd, err := cacher.Get(ctx, sku)
	assert.NoError(t, err)
	assert.Equal(t, product, retProd)

	cachedProd, err := cacher.Get(ctx, sku)
	assert.NoError(t, err)
	assert.Equal(t, product, cachedProd)

	nilProd, err := cacher.Get(ctx, skuNotExists)
	assert.ErrorIs(t, err, cart.ErrProductSkuNotFound)
	assert.Nil(t, nilProd)
}

func TestCacherRace(t *testing.T) {
	ctx := context.Background()
	provider := mock.NewProductProviderMock(t)

	totalProducts := 50
	products := make([]*model.Product, totalProducts)
	for i := 0; i < totalProducts; i++ {

		product := &model.Product{Name: fmt.Sprintf("Product %d", i), Price: model.Price(i * 10)}
		products[i] = product

		provider = provider.GetMock.Optional().When(ctx, model.Sku(i)).Then(product, nil)
	}

	inmemoryCache := inmemory.NewCache[CacheItem](CacheItem{})
	cacher := NewCacher(provider, inmemoryCache)
	cacher.removeSkuMxDuration = time.Second / 5
	cacher.cacheItemTTL = time.Second / 3

	var wg sync.WaitGroup

	n := 100
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			skuI := rand.Int63n(int64(totalProducts))

			sku := model.Sku(skuI)
			product := products[skuI]

			prod, err := cacher.Get(ctx, sku)

			assert.NoError(t, err)
			assert.Equal(t, product, prod)
		}()
	}

	wg.Wait()

	time.Sleep(time.Second / 2)

	mapLen := 0
	cacher.skuMxMap.Range(func(key, value any) bool {
		mapLen++
		return true
	})
	assert.Equal(t, 0, mapLen)
}
