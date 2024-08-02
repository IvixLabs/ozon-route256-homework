package e2e

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"os"
	appGrpc "route256/loms/internal/app/grpc"
	desc "route256/loms/internal/pb/loms/v1"
	"testing"
)

type GRPCSuite struct {
	suite.Suite
	ctx        context.Context
	container  testcontainers.Container
	host       string
	port       string
	grpcConn   *grpc.ClientConn
	grpcClient desc.LomsClient
}

func (s *GRPCSuite) SetupSuite() {
	s.ctx = context.Background()

	dockerContextPath := os.Getenv("DOCKER_CONTEXT_PATH")
	if dockerContextPath == "" {
		panic("DOCKER_CONTEXT_PATH is required")
	}

	const AppGrpcPort = "50051"

	reqContainers := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    dockerContextPath,
			Dockerfile: "build/Dockerfile",
			KeepImage:  false,
			BuildOptionsModifier: func(buildOptions *types.ImageBuildOptions) {
				buildOptions.Target = "run"
			},
		},
		ExposedPorts: []string{AppGrpcPort},
		Env: map[string]string{
			"APP_GRPC_ADDR":                 ":" + AppGrpcPort,
			"APP_STORAGE_MODE":              string(appGrpc.InMemoryStorageMode),
			"APP_KAFKA_MODE":                string(appGrpc.NoopKafkaMode),
			"APP_GRACEFUL_SHUTDOWN_TIMEOUT": "5",
			"APP_DEBUG_SRV":                 "disabled",
			"APP_TRACER":                    "disabled",
			"APP_GRPCGW":                    "disabled",
		},
		WaitingFor: wait.ForListeningPort(AppGrpcPort),
	}

	var err error
	s.container, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: reqContainers,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	cHost, err := s.container.Host(s.ctx)
	if err != nil {
		panic(err)
	}
	s.host = cHost

	cPort, err := s.container.MappedPort(s.ctx, AppGrpcPort)
	if err != nil {
		panic(err)
	}
	s.port = cPort.Port()

	var grpcErr error
	s.grpcConn, grpcErr = grpc.NewClient(s.host+":"+s.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if grpcErr != nil {
		panic(grpcErr)
	}

	s.grpcClient = desc.NewLomsClient(s.grpcConn)
}

func (s *GRPCSuite) TearDownSuite() {
	_ = s.grpcConn.Close()
	_ = s.container.Terminate(s.ctx)
}

func (s *GRPCSuite) SetupTest() {
}

func (s *GRPCSuite) TearDownTest() {
}
func (s *GRPCSuite) TestGrpc() {

	var createdOrderID int64

	tests := []struct {
		name      string
		operation func(*testing.T)
	}{
		{
			name: "Order_info_not_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderInfo(s.ctx, &desc.OrderInfoRequest{OrderID: 1})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "order info: order is not found", stat.Message())
			},
		},
		{
			name: "Order_cancel_not_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCancel(s.ctx, &desc.OrderCancelRequest{OrderID: 1})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, stat.Code())
				assert.Equal(t, "order cancel: order is not found", stat.Message())
			},
		},
		{
			name: "Order_pay_not_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderPay(s.ctx, &desc.OrderPayRequest{OrderID: 1})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "order pay: order is not found", stat.Message())
			},
		},
		{
			name: "Stock_info_not_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.StockInfo(s.ctx, &desc.StockInfoRequest{Sku: 123})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "stock info: stock is not found", stat.Message())
			},
		},
		{
			name: "Stock_info_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.StockInfo(s.ctx, &desc.StockInfoRequest{Sku: 1076963})
				assert.NoError(t, err)
				assert.Equal(t, uint64(4), res.Count)
			},
		},
		{
			name: "Order_creat_stock_not_found",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCreate(s.ctx,
					&desc.OrderCreateRequest{
						User: 1,
						Items: []*desc.OrderItem{
							{
								Sku:   1,
								Count: 1,
							},
						},
					})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "order create: stock is not found", stat.Message())
			},
		},
		{
			name: "Order_creat_insufficient_count",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCreate(s.ctx,
					&desc.OrderCreateRequest{
						User: 1,
						Items: []*desc.OrderItem{
							{
								Sku:   1076963,
								Count: 5,
							},
						},
					})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "order create: insufficient stock count", stat.Message())
			},
		},
		{
			name: "Order_creat_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCreate(s.ctx,
					&desc.OrderCreateRequest{
						User: 1,
						Items: []*desc.OrderItem{
							{
								Sku:   1076963,
								Count: 2,
							},
						},
					})
				assert.NoError(t, err)
				assert.NotEqual(t, int64(0), res.OrderID)

				createdOrderID = res.OrderID
			},
		},
		{
			name: "Order_info_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderInfo(s.ctx,
					&desc.OrderInfoRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.Equal(t, "awaiting payment", res.Status)
				assert.Equal(t, []*desc.OrderItem{{Sku: 1076963, Count: 2}}, res.Items)
			},
		},
		{
			name: "Order_add_new_items_insufficient_count",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCreate(s.ctx,
					&desc.OrderCreateRequest{
						User: 1,
						Items: []*desc.OrderItem{
							{
								Sku:   1076963,
								Count: 3,
							},
						},
					})
				assert.Nil(t, res)

				stat, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.FailedPrecondition, stat.Code())
				assert.Equal(t, "order create: insufficient stock count", stat.Message())
			},
		},
		{
			name: "Order_cancel_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCancel(s.ctx,
					&desc.OrderCancelRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "Order_info_cancelled_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderInfo(s.ctx,
					&desc.OrderInfoRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.Equal(t, "cancelled", res.Status)
				assert.Equal(t, []*desc.OrderItem{{Sku: 1076963, Count: 2}}, res.Items)
			},
		},
		{
			name: "Order_creat_ok_full_count",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderCreate(s.ctx,
					&desc.OrderCreateRequest{
						User: 1,
						Items: []*desc.OrderItem{
							{
								Sku:   1076963,
								Count: 4,
							},
						},
					})
				assert.NoError(t, err)
				assert.NotEqual(t, int64(0), res.OrderID)

				createdOrderID = res.OrderID
			},
		},
		{
			name: "Order_info_full_count_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderInfo(s.ctx,
					&desc.OrderInfoRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.Equal(t, "awaiting payment", res.Status)
				assert.Equal(t, []*desc.OrderItem{{Sku: 1076963, Count: 4}}, res.Items)
			},
		},
		{
			name: "Order_pay_ok",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderPay(s.ctx,
					&desc.OrderPayRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.NotNil(t, res)
			},
		},
		{
			name: "Order_info_full_count_payed",
			operation: func(t *testing.T) {
				res, err := s.grpcClient.OrderInfo(s.ctx,
					&desc.OrderInfoRequest{
						OrderID: createdOrderID,
					})
				assert.NoError(t, err)
				assert.Equal(t, "payed", res.Status)
				assert.Equal(t, []*desc.OrderItem{{Sku: 1076963, Count: 4}}, res.Items)
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.operation(t)
		})
	}
}
