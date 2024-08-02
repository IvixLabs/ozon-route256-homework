package loms

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	desc "route256/cart/internal/pb/loms/v1"
)

type GrpcClient struct {
	client desc.LomsClient
	conn   *grpc.ClientConn
}

func NewGrpcClient(grpcAddr string) *GrpcClient {

	grpcConn, err := grpc.NewClient(grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		panic(err)
	}
	client := desc.NewLomsClient(grpcConn)

	return &GrpcClient{client: client, conn: grpcConn}
}

func (c *GrpcClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
