package stock

import (
	"context"
	"errors"
	"route256/loms/internal/model"
)

var ErrReservedStockNotFound = errors.New("reserved stock is not found")

type ReservedStockRepository interface {
	Save(ctx context.Context, rStock *model.ReservedStock) error
	GetLocked(ctx context.Context, orderID model.OrderID, sku model.Sku) (*model.ReservedStock, error)
}
