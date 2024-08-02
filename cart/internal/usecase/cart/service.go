package cart

import (
	"context"
	"errors"
	"route256/cart/internal/model"
)

var (
	ErrCartNotFound     = errors.New("cart is not found")
	ErrCartItemNotFound = errors.New("cart item is not found")
)

type Repository interface {
	GetCart(ctx context.Context, userId model.UserID) (*model.Cart, error)
	UpdateCart(ctx context.Context, cart *model.Cart)
	RemoveCart(ctx context.Context, userId model.UserID)
	HasCart(ctx context.Context, userId model.UserID) bool
}

type Service struct {
	repository      Repository
	productProvider ProductProvider
	lomsClient      LOMSClient
	config          Config
}

func NewService(config Config, repository Repository, productProvider ProductProvider, stockProvider LOMSClient) *Service {
	return &Service{repository: repository, productProvider: productProvider, lomsClient: stockProvider, config: config}
}
