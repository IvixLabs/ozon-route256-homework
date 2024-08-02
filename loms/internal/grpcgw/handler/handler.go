package handler

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"route256/loms/internal/pb/loms/v1"
)

type GrpcGw struct {
	conn *grpc.ClientConn
	mux  http.Handler
}

func NewHandler(grpcAddr string) (*GrpcGw, error) {
	grpcgw := &GrpcGw{}

	clientConn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcgw.conn = clientConn

	gwmux := runtime.NewServeMux()

	err = loms.RegisterLomsHandler(context.Background(), gwmux, clientConn)
	if err != nil {
		return nil, err
	}

	grpcgw.mux = gwmux

	return grpcgw, nil
}

func (g *GrpcGw) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

func (g *GrpcGw) Close() error {
	err := g.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
