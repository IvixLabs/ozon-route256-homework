package cart

import (
	"context"
	"errors"
	"route256/cart/internal/model"
)

var (
	ErrStockNotFound = errors.New("stock is not found")
)

type LOMSClient interface {
	GetStockCount(ctx context.Context, sku model.Sku) (model.Count, error)
	CreateOrder(ctx context.Context, cart *model.Cart) (model.OrderID, error)
}
