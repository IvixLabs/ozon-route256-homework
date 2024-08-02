package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPool(ctx context.Context, dbString string, dbReplicaString string) (Pool, error) {
	pgxConn, err := pgx.Connect(ctx, dbString)
	if err != nil {
		return nil, err
	}

	orderStatusType, err := pgxConn.LoadType(ctx, "order_status")
	if err != nil {
		return nil, err
	}
	pgxConn.TypeMap().RegisterType(orderStatusType)

	arrayOrderStatusType, err := pgxConn.LoadType(ctx, "_order_status")
	if err != nil {
		return nil, err
	}
	err = pgxConn.Close(ctx)
	if err != nil {
		return nil, err
	}

	dbMasterPool, err := createPgxPool(ctx, dbString, orderStatusType, arrayOrderStatusType)
	if err != nil {
		return nil, err
	}

	if dbReplicaString != "" {

		dbReplicaPool, err := createPgxPool(ctx, dbReplicaString, orderStatusType, arrayOrderStatusType)
		if err != nil {
			return nil, err
		}

		return NewPool(dbMasterPool, dbReplicaPool), nil
	}

	return dbMasterPool, nil
}

func createPgxPool(ctx context.Context, dbString string, orderStatusType *pgtype.Type, arrayOrderStatusType *pgtype.Type) (*pgxpool.Pool, error) {
	dbMasterPoolConfig, err := pgxpool.ParseConfig(dbString)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbMasterPoolConfig)

	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	dbMasterPoolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.TypeMap().RegisterType(orderStatusType)
		conn.TypeMap().RegisterType(arrayOrderStatusType)
		return nil
	}
	dbPool, err = pgxpool.NewWithConfig(ctx, dbMasterPoolConfig)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
