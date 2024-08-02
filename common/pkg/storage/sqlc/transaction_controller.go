package sqlc

import (
	"context"
	"route256/common/pkg/manager/transaction"
)

type TransactionController interface {
	transaction.Controller
	GetDBTX(ctx context.Context) DBTX
}
