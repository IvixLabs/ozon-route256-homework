package reservedstock

import (
	"context"
	"fmt"
	"route256/loms/internal/model"
	"route256/loms/internal/storage/inmemory"
	"route256/loms/internal/usecase/stock"
)

type InMemoryRepository struct {
	storage *inmemory.Storage
}

func NewInMemoryRepository(storage *inmemory.Storage) *InMemoryRepository {
	return &InMemoryRepository{storage: storage}
}

func (r *InMemoryRepository) Save(_ context.Context, rStock *model.ReservedStock) error {
	r.storage.LockReservedStocks()
	defer r.storage.UnlockReservedStocks()

	r.storage.ReservedStocks[getPK(rStock.OrderID, rStock.Sku)] = *rStock

	return nil
}

func getPK(orderID model.OrderID, sku model.Sku) string {
	return fmt.Sprintf("%d_%d", orderID, sku)
}

func (r *InMemoryRepository) GetLocked(_ context.Context, orderID model.OrderID, sku model.Sku) (*model.ReservedStock, error) {
	r.storage.LockReservedStocks()
	defer r.storage.UnlockReservedStocks()

	rStock, ok := r.storage.ReservedStocks[getPK(orderID, sku)]

	if !ok {

		return nil, stock.ErrReservedStockNotFound

	}

	return &rStock, nil
}
