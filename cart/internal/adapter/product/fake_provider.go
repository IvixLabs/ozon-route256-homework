package product

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
)

type FakeProvider struct {
	config   Config
	products map[model.Sku]model.Product
}

func NewFakeProvider() *FakeProvider {
	return &FakeProvider{
		products: map[model.Sku]model.Product{
			111: {Name: "Product 111", Price: 100},
			222: {Name: "Product 222", Price: 200},
			333: {Name: "Product 333", Price: 300},
		},
	}
}

func (p *FakeProvider) Get(_ context.Context, sku model.Sku) (model.Product, error) {
	product, ok := p.products[sku]

	if ok {
		return product, nil
	}

	return cart.EmptyProduct, cart.ErrProductSkuNotFound
}
