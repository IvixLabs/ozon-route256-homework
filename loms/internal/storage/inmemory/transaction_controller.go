package inmemory

import (
	"context"
	transaction2 "route256/common/pkg/manager/transaction"
	"route256/loms/internal/model"
	"sync"
)

type Transaction struct {
	commitFn   func()
	rollbackFn func()
}

func NewTransaction(commitFn func(), rollbackFn func()) *Transaction {
	return &Transaction{commitFn: commitFn, rollbackFn: rollbackFn}
}

func (m *Transaction) Commit(_ context.Context) error {
	m.commitFn()
	return nil
}

func (m *Transaction) Rollback(_ context.Context) error {
	m.rollbackFn()
	return nil
}

type Controller struct {
	muTransaction sync.Mutex
	storage       *Storage
}

func NewController(storage *Storage) *Controller {
	return &Controller{storage: storage}
}

func (c *Controller) Begin(ctx context.Context) (transaction2.Transaction, context.Context, error) {
	c.muTransaction.Lock()

	c.storage.LockOrders()
	ordersState := make(map[model.OrderID]model.Order)
	for k, v := range c.storage.Orders {
		ordersState[k] = v
	}
	c.storage.UnlockOrders()

	c.storage.LockStocks()
	stocksState := make(map[model.Sku]model.Stock)
	for k, v := range c.storage.Stocks {
		stocksState[k] = v
	}
	c.storage.UnlockStocks()

	return NewTransaction(func() {
		c.muTransaction.Unlock()
	}, func() {
		c.storage.Orders = ordersState
		c.storage.Stocks = stocksState
		c.muTransaction.Unlock()
	}), ctx, nil
}
