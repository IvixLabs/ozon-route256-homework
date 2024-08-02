package cart

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"route256/cart/internal/model"
	usecaseCart "route256/cart/internal/usecase/cart"
	"sync"
)

var cartsCurrentAmountVec = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "app_cart",
		Name:      "carts_current_amount",
		Help:      "Current carts amount",
	},
	[]string{},
)

type InMemoryRepository struct {
	carts map[model.UserID]model.Cart
	mu    sync.Mutex
}

func (i *InMemoryRepository) UpdateCart(ctx context.Context, cart *model.Cart) {
	ctx, span := i.beginSpan(ctx, "UpdateCart")
	defer span.End()

	i.mu.Lock()
	defer i.mu.Unlock()

	i.carts[cart.UserId] = *cart

	l := len(i.carts)
	go func() {
		cartsCurrentAmountVec.WithLabelValues().Set(float64(l))
	}()
}

func (i *InMemoryRepository) RemoveCart(ctx context.Context, userId model.UserID) {
	ctx, span := i.beginSpan(ctx, "RemoveCart")
	defer span.End()

	i.mu.Lock()
	defer i.mu.Unlock()

	delete(i.carts, userId)

	l := len(i.carts)
	go func() {
		cartsCurrentAmountVec.WithLabelValues().Set(float64(l))
	}()
}

func (i *InMemoryRepository) HasCart(ctx context.Context, userId model.UserID) bool {
	ctx, span := i.beginSpan(ctx, "HasCart")
	defer span.End()

	i.mu.Lock()
	defer i.mu.Unlock()

	_, ok := i.carts[userId]

	return ok
}

func (i *InMemoryRepository) GetCart(ctx context.Context, userId model.UserID) (*model.Cart, error) {
	ctx, span := i.beginSpan(ctx, "GetCart")
	defer span.End()

	i.mu.Lock()
	defer i.mu.Unlock()

	cart, ok := i.carts[userId]
	if !ok {
		return nil, usecaseCart.ErrCartNotFound
	}

	return &cart, nil
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		carts: make(map[model.UserID]model.Cart),
	}
}
