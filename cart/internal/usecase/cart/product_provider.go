package cart

import (
	"context"
	"errors"
	"route256/cart/internal/model"
)

var (
	ErrProductSkuNotFound    = errors.New("product sku is not found")
	ErrProductServiceIsBusy  = errors.New("product service is busy")
	ErrProductServiceProblem = errors.New("product service has wrong answer")
)

var EmptyProduct = model.Product{}

type ProductProvider interface {
	Get(ctx context.Context, sku model.Sku) (model.Product, error)
}
