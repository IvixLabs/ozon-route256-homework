package redis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"route256/cart/internal/cache"
	cacheRedis "route256/cart/internal/cache/redis"
	"testing"
	"time"
)

type CacheSuite struct {
	suite.Suite
	ctx       context.Context
	redisAddr string
	redisCont *redis.RedisContainer
}

func (s *CacheSuite) SetupSuite() {
	s.ctx = context.Background()

	var err error
	s.redisCont, err = redis.Run(s.ctx, "redis:6")

	if err != nil {
		panic(err)
	}

	host, err := s.redisCont.Host(s.ctx)
	port, err := s.redisCont.MappedPort(s.ctx, "6379")

	s.redisAddr = host + ":" + port.Port()
}

func (s *CacheSuite) TearDownSuite() {
	_ = s.redisCont.Terminate(s.ctx)

}

func (s *CacheSuite) SetupTest() {

}

func (s *CacheSuite) TearDownTest() {

}

type testItem struct {
	A string
	B int
}

func (s *CacheSuite) TestCache() {
	t := s.T()

	ctx := context.Background()
	item := testItem{A: "aaa", B: 111}
	key := "abc"
	emptyItem := testItem{}

	redisCache := cacheRedis.NewCache[testItem](s.redisAddr, emptyItem)

	t.Run("Connect", func(t *testing.T) {
		err := redisCache.Connect(ctx)
		assert.NoError(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		err := redisCache.Set(ctx, key, item, 2*time.Second)
		assert.NoError(t, err)
	})

	t.Run("Get", func(t *testing.T) {
		foundItem, err := redisCache.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, foundItem, item)
	})

	t.Run("Delete", func(t *testing.T) {
		err := redisCache.Delete(ctx, key)
		assert.NoError(t, err)
	})
	t.Run("Get not found", func(t *testing.T) {
		notFoundItem, err := redisCache.Get(ctx, key)
		assert.ErrorIs(t, err, cache.ErrItemNotFound)
		assert.Equal(t, notFoundItem, emptyItem)
	})

	t.Run("Set with expiration", func(t *testing.T) {
		err := redisCache.Set(ctx, key, item, time.Second/2)
		assert.NoError(t, err)

	})

	time.Sleep(time.Second)

	t.Run("Get with expiration", func(t *testing.T) {
		notFoundItem, err := redisCache.Get(ctx, key)
		assert.ErrorIs(t, err, cache.ErrItemNotFound)
		assert.Equal(t, emptyItem, notFoundItem)
	})

}
