package grpc

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	commonTransaction "route256/common/pkg/manager/transaction"
	"route256/logger/pkg/interceptor"
	"route256/logger/pkg/logger"
	grpcLoms "route256/loms/internal/app/grpc/service/loms"
	interceptor2 "route256/loms/internal/interceptor"
	logger2 "route256/loms/internal/logger"
	"route256/loms/internal/manager/shard"
	"route256/loms/internal/manager/transaction"
	"route256/loms/internal/pb/loms/v1"
	"route256/loms/internal/repository/order"
	"route256/loms/internal/repository/reservedstock"
	"route256/loms/internal/repository/stock"
	"route256/loms/internal/storage/inmemory"
	"route256/loms/internal/storage/sqlc"
	usecaseOrder "route256/loms/internal/usecase/order"
	usecaseStock "route256/loms/internal/usecase/stock"
	"route256/metrics/pkg/tracer"
	"route256/transactionalbox/pkg/transactionalbox"
	"route256/transactionalbox/pkg/transactionalbox/broker/kafka"
	"route256/transactionalbox/pkg/transactionalbox/broker/noop"
	"route256/transactionalbox/pkg/transactionalbox/store"
	"sync"
)

type Server struct {
	config Config
}

func NewServer(config Config) *Server {
	return &Server{config: config}
}

func (server *Server) Run(ctx context.Context) {

	ctx = logger.ContextWithLogger(ctx, logger.InfoLevel, "app_loms")
	tracer.InitTracerProvider(ctx, "app_loms")

	conn, err := net.Listen("tcp", server.config.GRPCAddr)
	if err != nil {
		log.Panicln(err)
	}

	var stockRepo usecaseStock.Repository
	var orderRepo usecaseOrder.Repository
	var reservedStockRepo usecaseStock.ReservedStockRepository
	var msgPublisherStore transactionalbox.Store

	var tcs []commonTransaction.Controller

	if server.config.StorageMode == SqlcStorageMode {
		dbPool, poolErr := sqlc.GetPool(ctx, server.config.DBString, server.config.DBReplicaString)
		if poolErr != nil {
			log.Panicln(poolErr)
		}
		sqcTransactionController := sqlc.NewTransactionController(dbPool, 0)
		tcs = append(tcs, sqcTransactionController)

		dbPool1, poolErr1 := sqlc.GetPool(ctx, server.config.DBString1, server.config.DBReplicaString1)
		if poolErr1 != nil {
			log.Panicln(poolErr1)
		}
		sqcTransactionController1 := sqlc.NewTransactionController(dbPool1, 1)
		tcs = append(tcs, sqcTransactionController1)

		loggerObj := logger2.NewLogger()
		stockRepo = stock.NewSqlcRepository(sqcTransactionController, loggerObj)

		reservedStockRepo = reservedstock.NewSqlcRepository(sqcTransactionController, loggerObj)

		sqlcTxs := []*sqlc.TransactionController{sqcTransactionController, sqcTransactionController1}
		shardManager := shard.NewSqlcManager(sqlcTxs)
		orderRepo = order.NewSqlcRepository(shardManager, loggerObj)

		msgPublisherStore = store.NewSqlc(sqcTransactionController)
	} else if server.config.StorageMode == InMemoryStorageMode {
		inMemStorage := inmemory.NewStorage()
		inmemory.InitStockData(inMemStorage)

		inMemTransactionController := inmemory.NewController(inMemStorage)
		tcs = append(tcs, inMemTransactionController)

		stockRepo = stock.NewInMemoryRepository(inMemStorage)
		orderRepo = order.NewInMemoryRepository(inMemStorage)
		reservedStockRepo = reservedstock.NewInMemoryRepository(inMemStorage)

		msgPublisherStore = store.NewInMemory()
	} else {
		log.Panicln("Wrong storage mode")
	}

	txManager := transaction.NewManager(tcs)

	stockService := usecaseStock.NewService(stockRepo, reservedStockRepo)

	msgPublisher := transactionalbox.NewPublisher(msgPublisherStore)

	orderService := usecaseOrder.NewService(orderRepo, msgPublisher, stockService, txManager)

	lomsCtrl := grpcLoms.NewService(nil, orderService, stockRepo, stockService)

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.StatsHandler(interceptor.NewLogger(logger.FromContext(ctx))),
		grpc.ChainUnaryInterceptor(interceptor2.Logging),
		grpc.ChainUnaryInterceptor(interceptor2.RequestCounter),
		grpc.ChainUnaryInterceptor(interceptor2.Validate),
	)

	reflection.Register(grpcServer)

	loms.RegisterLomsServer(grpcServer, lomsCtrl)

	go func() {
		if err = grpcServer.Serve(conn); err != nil {
			log.Panicln(err)
		}
	}()

	var kafkaBroker transactionalbox.ProducerBroker

	if server.config.KafkaMode == KafkaKafkaMode {
		kafkaBroker = kafka.NewKafkaProducerBroker()
	} else {
		kafkaBroker = noop.NewProducerBroker()
	}

	producer := transactionalbox.NewProducer(kafkaBroker, msgPublisherStore)

	shutdownWg := sync.WaitGroup{}

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		producer.Run(ctx)
	}()

	logger.Infow(ctx, "Loms service is started")
	defer logger.Infow(ctx, "Loms service is stopped")

	<-ctx.Done()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		if shErr := tracer.Shutdown(context.Background()); shErr != nil {
			log.Panicln(shErr)
		}
	}()

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		grpcServer.GracefulStop()
	}()

	shutdownWg.Wait()
}
