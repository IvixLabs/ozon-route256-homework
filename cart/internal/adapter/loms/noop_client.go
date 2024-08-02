package loms

import (
	"context"
	"route256/cart/internal/model"
)

type NoopClient struct {
}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) CreateOrder(_ context.Context, _ *model.Cart) (model.OrderID, error) {
	return 1, nil
}

func (c *NoopClient) GetStockCount(_ context.Context, _ model.Sku) (model.Count, error) {
	return 100, nil
}
