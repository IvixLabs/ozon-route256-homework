package transaction

import (
	"context"
	"route256/common/pkg/manager/transaction"
)

type Manager struct {
	txControllers []transaction.Controller
}

func NewManager(txControllers []transaction.Controller) *Manager {
	return &Manager{
		txControllers: txControllers,
	}
}

type txItem struct {
	tx     transaction.Transaction
	closed bool
}

func rollbackTxItems(ctx context.Context, txItems []*txItem) {
	for _, tx := range txItems {
		if !tx.closed {
			commitErr := tx.tx.Rollback(ctx)
			if commitErr != nil {
				panic(commitErr)
			}
			tx.closed = true
		}
	}
}

func commitTxItems(ctx context.Context, txItems []*txItem) {
	for _, tx := range txItems {
		if !tx.closed {
			rollbackErr := tx.tx.Commit(ctx)
			if rollbackErr != nil {
				panic(rollbackErr)
			}
			tx.closed = true
		}
	}
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {

	newCtx := ctx

	txItems := make([]*txItem, 0, len(m.txControllers))
	defer func() {
		rollbackTxItems(ctx, txItems)
	}()

	for _, txc := range m.txControllers {
		tx, locCtx, err := txc.Begin(newCtx)
		if err != nil {
			return err
		}
		newCtx = locCtx
		txItems = append(txItems, &txItem{tx: tx})
	}

	err := fn(newCtx)

	if err != nil {
		return err
	}

	commitTxItems(ctx, txItems)

	return nil
}
