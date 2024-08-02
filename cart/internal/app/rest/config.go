package rest

import (
	"route256/cart/internal/adapter/product"
)

type AppMode string

var (
	AppModeProd AppMode = "prod"
)

type Config struct {
	Mode               AppMode
	Address            string
	ProductProvider    product.Config
	LolmsGrpcAddr      string
	ProductProviderRps int
	JaegerEndpoint     string
	RedisAddr          string
}
