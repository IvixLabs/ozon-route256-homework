package cart

import (
	"context"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/model"
	usecaseCart "route256/cart/internal/usecase/cart"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestInMemoryRepositoryTable(t *testing.T) {
	t.Parallel()
	var (
		unknownUserId = model.UserID(999)
		userId        = model.UserID(1)
		sku           = model.Sku(111)
		skuCount      = model.Count(5)

		secondUserId   = model.UserID(2)
		secondSku      = model.Sku(222)
		secondSkuCount = model.Count(3)
	)
	ctx := context.Background()

	createCart := func(userId model.UserID, sku model.Sku, count model.Count) model.Cart {
		cart := model.Cart{
			UserId: userId,
			Items: map[model.Sku]model.CartItem{
				sku: {Sku: sku, Count: count},
			},
		}

		return cart
	}

	createNotEmptyCarts := func() map[model.UserID]model.Cart {
		carts := make(map[model.UserID]model.Cart)
		carts[userId] = createCart(userId, sku, skuCount)
		carts[secondUserId] = createCart(secondUserId, secondSku, secondSkuCount)
		return carts
	}

	createSecondItemCarts := func() map[model.UserID]model.Cart {
		carts := make(map[model.UserID]model.Cart)
		carts[secondUserId] = createCart(secondUserId, secondSku, secondSkuCount)
		return carts
	}

	createEmptyCarts := func() map[model.UserID]model.Cart {
		return map[model.UserID]model.Cart{}
	}

	tests := []struct {
		name      string
		carts     map[model.UserID]model.Cart
		operation func(t *testing.T, r *InMemoryRepository)
	}{
		{
			name:  "Empty_repo",
			carts: createEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {

			},
		},
		{
			name:  "Get_unknown_cart",
			carts: createEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				cart, err := r.GetCart(ctx, unknownUserId)
				assert.Nil(t, cart)
				assert.ErrorIs(t, err, usecaseCart.ErrCartNotFound)
			},
		},
		{
			name:  "Has_unknown_cart",
			carts: createEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				ok := r.HasCart(ctx, unknownUserId)
				assert.False(t, ok)
			},
		},
		{
			name:  "Remove_unknown_cart",
			carts: createEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				r.RemoveCart(ctx, unknownUserId)
			},
		},
		{
			name:  "Update_cart",
			carts: createNotEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				firstCart := createCart(userId, sku, skuCount)
				secondCart := createCart(secondUserId, secondSku, secondSkuCount)

				r.UpdateCart(ctx, &firstCart)
				r.UpdateCart(ctx, &secondCart)
			},
		},
		{
			name:  "Has_known_cart",
			carts: createNotEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				ok := r.HasCart(ctx, userId)
				assert.True(t, ok)
			},
		},
		{
			name:  "Get_known_cart",
			carts: createNotEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				cart, err := r.GetCart(ctx, userId)
				expCart := createCart(userId, sku, skuCount)

				assert.Equal(t, &expCart, cart)
				assert.NoError(t, err)
			},
		},
		{
			name:  "Remove_known_cart",
			carts: createSecondItemCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				r.RemoveCart(ctx, userId)
			},
		},
		{
			name:  "Remove_second_known_cart",
			carts: createEmptyCarts(),
			operation: func(t *testing.T, r *InMemoryRepository) {
				r.RemoveCart(ctx, secondUserId)
			},
		},
	}

	repo := NewInMemoryRepository()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.operation(t, repo)

			assert.Equal(t, repo.carts, tt.carts)
		})
	}
}

func TestInMemoryRepositoryConcurrency(t *testing.T) {

	ctx := context.Background()

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Update_cart",
			test: func(t *testing.T) {
				repo := NewInMemoryRepository()
				repo.mu.Lock()

				mu := sync.Mutex{}
				mu.Lock()

				go func() {
					mu.Lock()
					cart := model.NewCart(111)
					repo.UpdateCart(ctx, cart)
				}()

				mu.Unlock()
				time.Sleep(time.Millisecond)

				assert.Equal(t, 0, len(repo.carts))
				time.Sleep(time.Millisecond)

				repo.mu.Unlock()
				time.Sleep(time.Millisecond)

				repo.mu.Lock()
				assert.Equal(t, 1, len(repo.carts))
				repo.mu.Unlock()
			},
		},
		{
			name: "Has_cart",
			test: func(t *testing.T) {
				repo := NewInMemoryRepository()
				cart := model.NewCart(111)
				repo.UpdateCart(ctx, cart)

				repo.mu.Lock()

				var step atomic.Int32

				mu := sync.Mutex{}
				mu.Lock()

				go func() {
					mu.Lock()
					ok := repo.HasCart(ctx, 111)
					step.Add(1)
					assert.True(t, ok)
				}()

				mu.Unlock()
				assert.Equal(t, int32(0), step.Load())
				time.Sleep(time.Millisecond)

				repo.mu.Unlock()
				time.Sleep(time.Millisecond)

				assert.Equal(t, int32(1), step.Load())
			},
		},
		{
			name: "Get_cart",
			test: func(t *testing.T) {
				repo := NewInMemoryRepository()
				cart := model.NewCart(111)
				repo.UpdateCart(ctx, cart)
				repo.mu.Lock()

				var step atomic.Int32
				mu := sync.Mutex{}
				mu.Lock()

				go func() {
					mu.Lock()
					cart, err := repo.GetCart(ctx, 111)
					step.Add(1)
					assert.NotNil(t, cart)
					assert.NoError(t, err)

				}()

				mu.Unlock()
				time.Sleep(time.Millisecond)

				assert.Equal(t, int32(0), step.Load())

				repo.mu.Unlock()
				time.Sleep(time.Millisecond)

				assert.Equal(t, int32(1), step.Load())
			},
		},
		{
			name: "Remove_cart",
			test: func(t *testing.T) {
				repo := NewInMemoryRepository()
				cart := model.NewCart(111)
				repo.UpdateCart(ctx, cart)
				repo.mu.Lock()

				mu := sync.Mutex{}
				mu.Lock()

				go func() {
					mu.Lock()
					repo.RemoveCart(ctx, 111)
				}()

				mu.Unlock()
				time.Sleep(time.Millisecond)

				assert.Equal(t, 1, len(repo.carts))

				repo.mu.Unlock()
				time.Sleep(time.Millisecond)

				repo.mu.Lock()
				assert.Equal(t, 0, len(repo.carts))
				repo.mu.Unlock()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.test(t)
		})
	}

}

func BenchmarkInMemoryRepositoryUpdateCart(b *testing.B) {
	repo := NewInMemoryRepository()

	userId := model.UserID(1)

	ctx := context.Background()
	cart := &model.Cart{
		UserId: userId,
		Items: map[model.Sku]model.CartItem{
			111: {Sku: 111, Count: 1},
		},
	}

	for i := 0; i < b.N; i++ {
		repo.UpdateCart(ctx, cart)
	}
}

func BenchmarkInMemoryRepositoryUpdateRemoveCart(b *testing.B) {
	repo := NewInMemoryRepository()

	userId := model.UserID(1)

	ctx := context.Background()
	cart := &model.Cart{
		UserId: userId,
		Items: map[model.Sku]model.CartItem{
			111: {Sku: 111, Count: 1},
		},
	}

	for i := 0; i < b.N; i++ {
		repo.UpdateCart(ctx, cart)
		repo.RemoveCart(ctx, userId)
	}
}
