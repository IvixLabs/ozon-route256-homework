package order

import (
	"context"
	"errors"
	"route256/loms/internal/model"
)

var (
	ErrOrderNotFound = errors.New("order is not found")
)

type Repository interface {
	GetByID(ctx context.Context, orderID model.OrderID) (*model.Order, error)
	GetLockByID(ctx context.Context, orderID model.OrderID) (*model.Order, error)
	Save(ctx context.Context, order *model.Order) (*model.Order, error)
}
