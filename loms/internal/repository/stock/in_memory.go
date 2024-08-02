package stock

import (
	"context"
	_ "embed"
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

func (r *InMemoryRepository) GetBySku(ctx context.Context, sku model.Sku) (*model.Stock, error) {
	ctx, span := beginSpan(ctx, "inmemory", "GetBySku")
	defer span.End()

	r.storage.LockStocks()
	defer r.storage.UnlockStocks()

	stockObj, ok := r.storage.Stocks[sku]

	if !ok {
		return nil, stock.ErrStockNotFound
	}

	return &stockObj, nil
}

func (r *InMemoryRepository) GetLockBySku(ctx context.Context, sku model.Sku) (*model.Stock, error) {
	return r.GetBySku(ctx, sku)
}

func (r *InMemoryRepository) Save(ctx context.Context, stockObj *model.Stock) error {
	ctx, span := beginSpan(ctx, "inmemory", "Save")
	defer span.End()

	r.storage.LockStocks()
	defer r.storage.UnlockStocks()

	r.storage.Stocks[stockObj.Sku] = *stockObj

	return nil
}
