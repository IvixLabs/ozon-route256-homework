package stock

import (
	"context"
	"github.com/stretchr/testify/assert"
	"route256/loms/internal/model"
	"route256/loms/internal/storage/inmemory"
	"route256/loms/internal/usecase/stock"
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
			name: "Get_by_sku_error",
			operation: func(t *testing.T, r *InMemoryRepository) {
				stockObj, err := r.GetBySku(ctx, 111)
				assert.ErrorIs(t, err, stock.ErrStockNotFound)
				assert.Nil(t, stockObj)
			},
		},
		{
			name: "Save_ok",
			operation: func(t *testing.T, r *InMemoryRepository) {
				newStock := model.NewStock(222, 4)
				err := r.Save(ctx, newStock)

				assert.NoError(t, err)

				expectStock := *model.NewStock(222, 4)
				assert.Equal(t, expectStock, r.storage.Stocks[222])
			},
		},
		{
			name: "Get_by_sku_ok",
			operation: func(t *testing.T, r *InMemoryRepository) {
				stockObj, err := r.GetBySku(ctx, 222)
				assert.NoError(t, err)

				expectStock := model.NewStock(222, 4)
				assert.Equal(t, expectStock, stockObj)
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

	repo.storage.LockStocks()

	mu := sync.Mutex{}
	mu.Lock()

	go func() {
		mu.Lock()
		err := repo.Save(ctx, model.NewStock(model.Sku(123), model.Count(11)))
		assert.NoError(t, err)
	}()

	mu.Unlock()
	time.Sleep(time.Millisecond)

	assert.Equal(t, 0, len(repo.storage.Stocks))

	repo.storage.UnlockStocks()
	time.Sleep(time.Millisecond)

	repo.storage.LockStocks()
	assert.Equal(t, 1, len(repo.storage.Stocks))
	repo.storage.UnlockStocks()
}

func TestInMemoryRepository_Concurrency_get_by_sku(t *testing.T) {
	t.Parallel()

	storage := inmemory.NewStorage()
	repo := NewInMemoryRepository(storage)
	ctx := context.Background()

	err := repo.Save(ctx, model.NewStock(123, 11))
	assert.NoError(t, err)

	repo.storage.LockStocks()

	mu := sync.Mutex{}
	mu.Lock()

	step := atomic.Int32{}

	go func() {
		foundStock, err := repo.GetBySku(ctx, 123)
		step.Add(1)
		assert.NotNil(t, foundStock)
		assert.NoError(t, err)
	}()

	mu.Unlock()
	time.Sleep(time.Millisecond)

	assert.Equal(t, int32(0), step.Load())

	repo.storage.UnlockStocks()
	time.Sleep(time.Millisecond)

	assert.Equal(t, int32(1), step.Load())
}
