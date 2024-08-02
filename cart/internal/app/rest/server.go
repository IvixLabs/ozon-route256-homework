package rest

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"route256/cart/internal/adapter/loms"
	"route256/cart/internal/adapter/product"
	"route256/cart/internal/app/rest/handler"
	"route256/cart/internal/cache/inmemory"
	"route256/cart/internal/cache/redis"
	cartRepository "route256/cart/internal/repository/cart"
	"route256/cart/internal/usecase/cart"
	"route256/logger/pkg/logger"
	"route256/metrics/pkg/tracer"
	"sync"
)

type Server struct {
	config Config
}

func NewServer(config Config) *Server {
	return &Server{config: config}
}

func (s *Server) Run(ctx context.Context) {

	ctx = logger.ContextWithLogger(ctx, logger.InfoLevel, "app_cart")

	tracer.InitTracerProvider(ctx, "app_cart")

	var productProvider cart.ProductProvider
	var lomsClient cart.LOMSClient

	if s.config.Mode == AppModeProd {
		httpProvider := product.NewHTTPProvider(s.config.ProductProvider)

		redisCache := redis.NewCache[inmemory.NextCacheItem[product.CacheItem]](
			s.config.RedisAddr,
			inmemory.NextCacheItem[product.CacheItem]{},
		)
		redErr := redisCache.Connect(ctx)
		if redErr != nil {
			log.Panicln(redErr)
		}

		fastCache := inmemory.NewCache[product.CacheItem](product.CacheItem{})
		fastCache.SetNext(redisCache)

		productProvider = product.NewCacher(httpProvider, fastCache)

		grpcClient := loms.NewGrpcClient(s.config.LolmsGrpcAddr)
		defer func(grpcClient *loms.GrpcClient) {
			closeErr := grpcClient.Close()
			if closeErr != nil {
				log.Panicln(closeErr)
			}
		}(grpcClient)

		lomsClient = grpcClient
	} else {
		productProvider = product.NewFakeProvider()
		lomsClient = loms.NewNoopClient()
	}

	cartRepo := cartRepository.NewInMemoryRepository()
	cartServiceConfig := cart.Config{
		ProductProviderRps: s.config.ProductProviderRps,
	}
	cartService := cart.NewService(cartServiceConfig, cartRepo, productProvider, lomsClient)

	router := handler.NewRouter(cartService)

	httpConn, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		log.Panic(err)
	}

	mux := router.GetMux()

	httpServer := &http.Server{
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		if err := httpServer.Serve(httpConn); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Panicln(err)
			}
		}
	}()

	logger.Infow(ctx, "Cart service is started", "mode", s.config.Mode)
	defer logger.Infow(ctx, "Cart service is stopped")

	<-ctx.Done()

	shutdownWg := sync.WaitGroup{}

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Panicln(err)
		}
	}()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if err := tracer.Shutdown(context.Background()); err != nil {
			log.Panicln(err)
		}
	}()

	shutdownWg.Wait()
}
