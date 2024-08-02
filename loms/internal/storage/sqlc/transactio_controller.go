package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
	"route256/common/pkg/manager/transaction"
	"route256/common/pkg/storage/sqlc"
)

type TCKey int

type ctxKey struct {
	key TCKey
}

type TransactionController struct {
	pool Pool
	key  TCKey
}

func NewTransactionController(pool Pool, key TCKey) *TransactionController {
	return &TransactionController{pool: pool, key: key}
}

func (tc *TransactionController) Begin(ctx context.Context) (transaction.Transaction, context.Context, error) {
	tx, err := tc.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	tr := NewTransaction(tx)

	return tr, context.WithValue(ctx, ctxKey{key: tc.key}, tx), nil
}

func (tc *TransactionController) GetDBTX(ctx context.Context) sqlc.DBTX {
	tx := ctx.Value(ctxKey{key: tc.key})
	if tx != nil {
		return tx.(pgx.Tx)
	}

	return tc.pool
}
