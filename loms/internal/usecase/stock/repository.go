package stock

import (
	"context"
	"errors"
	"route256/loms/internal/model"
)

var (
	ErrStockNotFound = errors.New("stock is not found")
)

type Repository interface {
	GetBySku(ctx context.Context, sku model.Sku) (*model.Stock, error)
	GetLockBySku(ctx context.Context, sku model.Sku) (*model.Stock, error)
	Save(ctx context.Context, stock *model.Stock) error
}
