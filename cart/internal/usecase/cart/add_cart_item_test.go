package cart

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart/mock"
	"testing"
)

func TestAddCartItemTable(t *testing.T) {
	t.Parallel()

	var (
		userId    = model.UserID(1)
		firstSku  = model.Sku(111)
		secondSku = model.Sku(222)
	)

	createOneItemCart := func() *model.Cart {
		oneItemCart := model.Cart{
			UserId: userId,
			Items: map[model.Sku]model.CartItem{
				firstSku: {Sku: firstSku, Count: 2},
			},
		}

		return &oneItemCart
	}

	createTwoItemsCart := func() *model.Cart {
		twoItemsCart := model.Cart{
			UserId: userId,
			Items: map[model.Sku]model.CartItem{
				firstSku:  {Sku: firstSku, Count: 2},
				secondSku: {Sku: secondSku, Count: 1},
			}}

		return &twoItemsCart
	}

	createFirstProduct := func() model.Product {
		product := model.Product{Name: "Product 111", Price: 100}
		return product
	}

	ctrl := minimock.NewController(t)
	productProvider := mock.NewProductProviderMock(ctrl)
	cartRepository := mock.NewRepositoryMock(ctrl)
	lomsClient := mock.NewLOMSClientMock(ctrl)
	service := NewService(Config{ProductProviderRps: 10}, cartRepository, productProvider, lomsClient)

	ctx := context.Background()

	errUnexpected := errors.New("unexpected error")

	tests := []struct {
		name       string
		wantErr    error
		wantCart   func() *model.Cart
		cart       func() *model.Cart
		cartErr    error
		count      model.Count
		product    model.Product
		sku        model.Sku
		productErr error
		stockCount model.Count
		stockErr   error
	}{
		{
			name:    "Unexpected_error",
			count:   model.Count(1),
			wantErr: errUnexpected,
			wantCart: func() *model.Cart {
				return nil
			},
			cart: func() *model.Cart {
				return nil
			},
			cartErr: errUnexpected,
		},
		{
			name:    "Wrong count",
			count:   model.Count(0),
			wantErr: ErrWrongCartItemCount,
			wantCart: func() *model.Cart {
				return nil
			},
			cart: func() *model.Cart {
				return nil
			},
		},
		{
			name:    "Has_not_cart",
			product: createFirstProduct(),
			sku:     firstSku,
			count:   model.Count(2),
			wantCart: func() *model.Cart {
				return createOneItemCart()
			},
			cart: func() *model.Cart {
				return nil
			},
			cartErr:    ErrCartNotFound,
			stockCount: 2,
		},
		{
			name: "Has_cart_not_sku",
			cart: func() *model.Cart {
				return createOneItemCart()
			},
			product: createFirstProduct(),
			sku:     secondSku,
			count:   model.Count(1),
			wantCart: func() *model.Cart {
				return createTwoItemsCart()
			},
			stockCount: 2,
		},
		{
			name: "Has_cart_and_sku",
			cart: func() *model.Cart {
				return createOneItemCart()
			},
			sku:     firstSku,
			product: createFirstProduct(),
			count:   model.Count(2),
			wantCart: func() *model.Cart {
				cart := createOneItemCart()
				cartItem := cart.Items[firstSku]
				cartItem.Count = 4
				cart.Items[firstSku] = cartItem
				return cart
			},
			stockCount: 4,
		},
		{
			name:       "Wrong_sku",
			product:    EmptyProduct,
			productErr: ErrProductSkuNotFound,
			count:      model.Count(2),
			wantCart: func() *model.Cart {
				return nil
			},
			cart: func() *model.Cart {
				return nil
			},
			cartErr: ErrCartNotFound,
			wantErr: ErrProductSkuNotFound,
		},
		{
			name: "Has_cart_not_sku_stock_error",
			cart: func() *model.Cart {
				return createOneItemCart()
			},
			product: createFirstProduct(),
			sku:     secondSku,
			count:   model.Count(1),
			wantCart: func() *model.Cart {
				return nil
			},
			stockErr: errUnexpected,
			wantErr:  errUnexpected,
		},
		{
			name: "Has_cart_and_sku_low_stock",
			cart: func() *model.Cart {
				return createOneItemCart()
			},
			sku:     firstSku,
			product: createFirstProduct(),
			count:   model.Count(1),
			wantCart: func() *model.Cart {
				cart := createOneItemCart()
				cartItem := cart.Items[firstSku]
				cartItem.Count = 3
				cart.Items[firstSku] = cartItem
				return cart
			},
			stockCount: 2,
			wantErr:    ErrInsufficientStocks,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wantCart := tt.wantCart()
			cart := tt.cart()

			lomsClient.GetStockCountMock.Expect(ctx, tt.sku).Return(tt.stockCount, tt.stockErr)
			productProvider.GetMock.Expect(ctx, tt.sku).Return(tt.product, tt.productErr)
			cartRepository.GetCartMock.Expect(ctx, userId).Return(cart, tt.cartErr)
			cartRepository.UpdateCartMock.Expect(ctx, wantCart)

			err := service.AddCartItem(ctx, userId, tt.sku, tt.count)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
