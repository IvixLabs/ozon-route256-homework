package grpcgw

import (
	"context"
	"os"
	"route256/common/pkg/env"
)

const (
	EnvAppGRPCGW         = "APP_GRPCGW"
	EnvAppGrpcAddr       = "APP_GRPCGW_GRPC_ADDR"
	EnvAppGRPCGWRESTAddr = "APP_GRPCGW_REST_ADDR"
)

func Run(ctx context.Context) {
	if !isEnvGRPCGW() {
		return
	}

	grpcgwConfig := getConfig()
	grpcgwServer := NewServer(grpcgwConfig)
	grpcgwServer.Run(ctx)
}

func getConfig() Config {
	return Config{
		HTTPAddr: env.GetEnvVar(EnvAppGRPCGWRESTAddr),
		GRPCAddr: env.GetEnvVar(EnvAppGrpcAddr),
	}

}

func isEnvGRPCGW() bool {
	v, ok := os.LookupEnv(EnvAppGRPCGW)
	if !ok {
		return true
	}

	return v != "disabled"
}
