package order

import (
	"context"
	"route256/loms/internal/model"
	"route256/loms/internal/storage/inmemory"
	usecaseOrder "route256/loms/internal/usecase/order"
)

type InMemoryRepository struct {
	storage *inmemory.Storage
	lastId  uint64
}

func NewInMemoryRepository(storage *inmemory.Storage) *InMemoryRepository {
	return &InMemoryRepository{storage: storage}
}

func (r *InMemoryRepository) Save(ctx context.Context, order *model.Order) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "inmemory", "Save")
	defer span.End()

	r.storage.LockOrders()
	defer r.storage.UnlockOrders()

	newOrder := *order

	if newOrder.ID == 0 {
		r.lastId++
		newOrder.ID = model.OrderID(r.lastId)
	}

	r.storage.Orders[newOrder.ID] = newOrder

	return &newOrder, nil
}

func (r *InMemoryRepository) GetByID(ctx context.Context, orderId model.OrderID) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "inmemory", "GetByID")
	defer span.End()

	r.storage.LockOrders()
	defer r.storage.UnlockOrders()

	order, ok := r.storage.Orders[orderId]

	if !ok {
		return nil, usecaseOrder.ErrOrderNotFound
	}

	return &order, nil
}

func (r *InMemoryRepository) GetLockByID(ctx context.Context, orderId model.OrderID) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "inmemory", "GetLockByID")
	defer span.End()

	return r.GetByID(ctx, orderId)
}
