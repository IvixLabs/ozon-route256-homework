package order

import (
	"context"
	"github.com/stretchr/testify/assert"
	"route256/loms/internal/model"
	"route256/loms/internal/storage/inmemory"
	"route256/loms/internal/usecase/order"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestInMemoryRepositoryTable(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewStorage()
	repo := NewInMemoryRepository(storage)
	ctx := context.Background()

	tests := []struct {
		name      string
		operation func(t *testing.T, r *InMemoryRepository)
	}{
		{
			name: "Get_by_id_error",
			operation: func(t *testing.T, r *InMemoryRepository) {
				stock, err := r.GetByID(ctx, 111)
				assert.ErrorIs(t, err, order.ErrOrderNotFound)
				assert.Nil(t, stock)
			},
		},
		{
			name: "Save_first",
			operation: func(t *testing.T, r *InMemoryRepository) {
				orderObj, err := r.Save(ctx, &model.Order{})
				assert.Equal(t, &model.Order{ID: 1}, orderObj)
				assert.Equal(t, uint64(1), r.lastId)
				assert.NoError(t, err)
			},
		},
		{
			name: "Save_second",
			operation: func(t *testing.T, r *InMemoryRepository) {
				orderObj, err := r.Save(ctx, &model.Order{})
				assert.Equal(t, &model.Order{ID: 2}, orderObj)
				assert.Equal(t, uint64(2), r.lastId)
				assert.NoError(t, err)
			},
		},
		{
			name: "Get_by_id_ok",
			operation: func(t *testing.T, r *InMemoryRepository) {
				stock, err := r.GetByID(ctx, 1)
				assert.NoError(t, err)
				assert.Equal(t, stock, &model.Order{ID: model.OrderID(1)})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.operation(t, repo)
		})
	}

}

func TestInMemoryRepository_Concurrency_save(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewStorage()
	repo := NewInMemoryRepository(storage)
	ctx := context.Background()

	repo.storage.LockOrders()

	mu := sync.Mutex{}
	mu.Lock()

	go func() {
		mu.Lock()
		newOrder, err := repo.Save(ctx, model.NewOrder(123, nil))
		assert.NotNil(t, newOrder)
		assert.NoError(t, err)
	}()

	mu.Unlock()
	time.Sleep(time.Millisecond)

	assert.Equal(t, 0, len(repo.storage.Orders))

	repo.storage.UnlockOrders()
	time.Sleep(time.Millisecond)

	repo.storage.LockOrders()
	assert.Equal(t, 1, len(repo.storage.Orders))
	repo.storage.UnlockOrders()
}

func TestInMemoryRepository_Concurrency_get_by_id(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewStorage()
	repo := NewInMemoryRepository(storage)
	ctx := context.Background()

	orderObj, _ := repo.Save(ctx, model.NewOrder(123, nil))

	repo.storage.LockOrders()

	mu := sync.Mutex{}
	mu.Lock()

	step := atomic.Int32{}

	go func() {
		mu.Lock()
		foundOrder, err := repo.GetByID(ctx, orderObj.ID)
		assert.NotNil(t, foundOrder)
		assert.NoError(t, err)
		step.Add(1)
	}()

	mu.Unlock()
	time.Sleep(time.Millisecond)

	assert.Equal(t, int32(0), step.Load())

	repo.storage.UnlockOrders()
	time.Sleep(time.Millisecond)

	assert.Equal(t, int32(1), step.Load())
}
