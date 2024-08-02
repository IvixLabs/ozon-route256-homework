package product

import (
	"context"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"route256/cart/internal/cache"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
	"strconv"
	"sync"
	"time"
)

type CacheItem struct {
	Product    model.Product
	HasProduct bool
}

type Cacher struct {
	storage             cache.Cache[CacheItem]
	provider            cart.ProductProvider
	skuMxMap            sync.Map
	removeSkuMxDuration time.Duration
	cacheItemTTL        time.Duration
}

func NewCacher(provider cart.ProductProvider, storage cache.Cache[CacheItem]) *Cacher {
	return &Cacher{
		storage:             storage,
		provider:            provider,
		removeSkuMxDuration: time.Second,
		cacheItemTTL:        20 * time.Second,
	}
}

var cacheHitCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app_cart",
		Name:      "product_cache_hit_total_counter",
		Help:      "Total amount of product cache hits",
	},
	[]string{},
)

var cacheMissCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app_cart",
		Name:      "product_cache_miss_total_counter",
		Help:      "Total amount of product cache misses",
	},
	[]string{},
)

var cacheHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app_cart",
		Name:      "product_cache_duration_histogram",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"type"})

func (c *Cacher) Get(ctx context.Context, sku model.Sku) (model.Product, error) {
	var cachedRes bool

	defer func(createdAt time.Time) {
		resType := "http"
		if cachedRes {
			resType = "cached"
		}

		cacheHistogramVec.WithLabelValues(resType).Observe(time.Since(createdAt).Seconds())
	}(time.Now())

	strSku := strconv.FormatInt(int64(sku), 10)

	cItem, err := c.storage.Get(ctx, strSku)

	if err != nil {
		if !errors.Is(err, cache.ErrItemNotFound) {
			return cart.EmptyProduct, err
		}
	} else {
		cacheHitCounterVec.WithLabelValues().Inc()
		cachedRes = true
		return expandCacheItem(cItem)
	}

	var skuMx sync.Locker = &sync.Mutex{}
	skuMx.Lock()
	defer skuMx.Unlock()

	actualSkuMx, skuMxOk := c.skuMxMap.LoadOrStore(sku, skuMx)

	foundSkuMx := actualSkuMx.(sync.Locker)
	if foundSkuMx != skuMx {
		foundSkuMx.Lock()
		defer foundSkuMx.Unlock()
	}

	if skuMxOk {
		cItem, err = c.storage.Get(ctx, strSku)

		if err != nil {
			return cart.EmptyProduct, err
		}

		cacheHitCounterVec.WithLabelValues().Inc()
		cachedRes = true
		return expandCacheItem(cItem)
	} else {
		defer func() {
			go func() {
				time.Sleep(c.removeSkuMxDuration)
				c.skuMxMap.Delete(sku)
			}()
		}()
	}

	product, err := c.provider.Get(ctx, sku)

	cacheItemObj := CacheItem{}
	if err != nil {
		if !errors.Is(err, cart.ErrProductSkuNotFound) {
			return cart.EmptyProduct, err
		}
	}

	if product != cart.EmptyProduct {
		cacheItemObj.Product = product
		cacheItemObj.HasProduct = true
	}

	err = c.storage.Set(ctx, strSku, cacheItemObj, c.cacheItemTTL)
	if err != nil {
		return cart.EmptyProduct, err
	}

	cacheMissCounterVec.WithLabelValues().Inc()
	return expandCacheItem(cacheItemObj)
}

func expandCacheItem(cItem CacheItem) (model.Product, error) {
	if !cItem.HasProduct {
		return cart.EmptyProduct, cart.ErrProductSkuNotFound
	}

	return cItem.Product, nil
}
