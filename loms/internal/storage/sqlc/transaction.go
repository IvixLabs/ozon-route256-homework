package sqlc

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	tx pgx.Tx
}

func NewTransaction(tx pgx.Tx) *Transaction {
	return &Transaction{
		tx: tx,
	}
}

func (tran *Transaction) Commit(ctx context.Context) error {
	err := tran.tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (tran *Transaction) Rollback(ctx context.Context) error {
	err := tran.tx.Rollback(ctx)
	if err != nil {
		return err
	}

	return nil
}
