package main

import (
	"context"
	"log"
	"os"
	"route256/cart/internal/adapter/product"
	"route256/cart/internal/app/rest"
	"route256/common/pkg/env"
	"strconv"
)

const (
	EnvAppMode                 = "APP_MODE"
	EnvAppRestAddr             = "APP_REST_ADDR"
	EnvAppProductProviderUrl   = "APP_PRODUCT_PROVIDER_URL"
	EnvAppProductProviderToken = "APP_PRODUCT_PROVIDER_TOKEN"
	EnvAppLomsGrpcAddr         = "APP_LOMS_GRPC_ADDR"
	EnvAppProductProviderRps   = "APP_PRODUCT_PROVIDER_RPS"
	EnvAppRedisAddr            = "APP_REDIS_ADDR"
)

func runRest(ctx context.Context) {
	restConfig := getRestConfig()
	restServer := rest.NewServer(restConfig)
	restServer.Run(ctx)
}

func getRestConfig() rest.Config {

	appMode := rest.AppMode(env.GetEnvVar(EnvAppMode))

	var appProductProviderUrl string
	var appProductProviderToken string
	var appLomsGrpcAddr string
	var appRedisAddr string

	if appMode == rest.AppModeProd {
		appProductProviderUrl = env.GetEnvVar(EnvAppProductProviderUrl)
		appProductProviderToken = env.GetEnvVar(EnvAppProductProviderToken)
		appLomsGrpcAddr = env.GetEnvVar(EnvAppLomsGrpcAddr)
		appRedisAddr = env.GetEnvVar(EnvAppRedisAddr)
	}

	appAddress := env.GetEnvVar(EnvAppRestAddr)

	strProductProviderRps := findEnvVar(EnvAppProductProviderRps)
	if strProductProviderRps == "" {
		strProductProviderRps = "10"
	}

	productProviderRps, err := strconv.Atoi(strProductProviderRps)
	if err != nil {
		log.Panicln(err)
	}

	return rest.Config{
		Mode:          appMode,
		Address:       appAddress,
		LolmsGrpcAddr: appLomsGrpcAddr,
		ProductProvider: product.Config{
			Url:   appProductProviderUrl,
			Token: appProductProviderToken,
		},
		ProductProviderRps: productProviderRps,
		RedisAddr:          appRedisAddr,
	}
}

func findEnvVar(name string) string {
	return os.Getenv(name)
}
