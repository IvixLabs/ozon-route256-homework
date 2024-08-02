package integration

import (
	"github.com/stretchr/testify/suite"
	"route256/cart/test/suite/redis"
	"testing"
)

func TestCache(t *testing.T) {
	suite.Run(t, &redis.CacheSuite{})
}
