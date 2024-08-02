package redis

import (
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"route256/cart/internal/cache"
	"time"
)

type Cache[V any] struct {
	addr      string
	client    *redis.Client
	emptyItem V
}

func NewCache[V any](addr string, emptyItem V) *Cache[V] {
	return &Cache[V]{addr: addr, emptyItem: emptyItem}
}

func (r *Cache[V]) Connect(ctx context.Context) error {
	r.client = redis.NewClient(&redis.Options{
		Addr:     r.addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := r.client.Ping(ctx)
	err := status.Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Cache[V]) Set(ctx context.Context, k string, v V, expiration time.Duration) error {

	strVal, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, k, strVal, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Cache[V]) Get(ctx context.Context, k string) (V, error) {
	val, err := r.client.Get(ctx, k).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return r.emptyItem, cache.ErrItemNotFound
		}

		return r.emptyItem, err
	}

	bufItem := r.emptyItem
	pItem := &bufItem
	err = json.Unmarshal([]byte(val), pItem)
	if err != nil {
		return r.emptyItem, err
	}

	return bufItem, nil
}

func (r *Cache[V]) Delete(ctx context.Context, k string) error {
	err := r.client.Del(ctx, k).Err()
	if err != nil {
		return err
	}

	return nil
}
