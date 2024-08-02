package main

import (
	"context"
	"route256/common/pkg/env"
	"route256/loms/internal/app/grpc"
)

const (
	EnvAppGrpcAddr       = "APP_GRPC_ADDR"
	EnvPgDbSting         = "APP_PG_DB_STRING"
	EnvPgDbReplicaSting  = "APP_PG_DB_REPLICA_STRING"
	EnvStorageMode       = "APP_STORAGE_MODE"
	EnvKafkaMode         = "APP_KAFKA_MODE"
	EnvPgDbSting1        = "APP_PG_DB_STRING1"
	EnvPgDbReplicaSting1 = "APP_PG_DB_REPLICA_STRING1"
)

func runGRPC(ctx context.Context) {

	grpcConfig := getGRPCConfig()
	grpcServer := grpc.NewServer(grpcConfig)
	grpcServer.Run(ctx)
}

func getGRPCConfig() grpc.Config {
	storageMode := grpc.StorageMode(env.GetEnvVar(EnvStorageMode))
	kafkaMode := grpc.KafkaMode(env.GetEnvVar(EnvKafkaMode))

	var dbString, dbReplicaString, dbString1, dbReplicaString1 string
	if storageMode == grpc.SqlcStorageMode {
		dbString = env.GetEnvVar(EnvPgDbSting)
		dbReplicaString = env.GetEnvVar(EnvPgDbReplicaSting)
		dbString1 = env.GetEnvVar(EnvPgDbSting1)
		dbReplicaString1 = env.GetEnvVar(EnvPgDbReplicaSting1)
	}

	return grpc.Config{
		GRPCAddr:         env.GetEnvVar(EnvAppGrpcAddr),
		DBString:         dbString,
		DBReplicaString:  dbReplicaString,
		StorageMode:      storageMode,
		KafkaMode:        kafkaMode,
		DBString1:        dbString1,
		DBReplicaString1: dbReplicaString1,
	}
}
