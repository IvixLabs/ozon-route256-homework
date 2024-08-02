package cart

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart/mock"
	"testing"
)

func TestRemoveCartItemTable(t *testing.T) {
	t.Parallel()

	sku := model.Sku(111)
	oneItemCart := &model.Cart{Items: map[model.Sku]model.CartItem{
		sku: {Sku: sku, Count: 1},
	}}

	secondSku := model.Sku(222)
	twoItemsCart := &model.Cart{Items: map[model.Sku]model.CartItem{
		sku:       {Sku: sku, Count: 1},
		secondSku: {Sku: secondSku, Count: 1},
	}}

	twoItemsCartUpdated := &model.Cart{Items: map[model.Sku]model.CartItem{
		secondSku: {Sku: secondSku, Count: 1},
	}}

	tests := []struct {
		name       string
		userId     model.UserID
		want       error
		cartErr    error
		cart       *model.Cart
		updateCart *model.Cart
	}{
		{
			name:    "Has_not_cart",
			userId:  model.UserID(1),
			want:    ErrCartNotFound,
			cartErr: ErrCartNotFound,
		},
		{
			name:   "Has_cart_not_sku",
			userId: model.UserID(1),
			want:   ErrCartItemNotFound,
			cart:   &model.Cart{Items: map[model.Sku]model.CartItem{}},
		},
		{
			name:   "Has_cart_and_sku",
			userId: model.UserID(1),
			want:   nil,
			cart:   oneItemCart,
		},
		{
			name:       "Has_cart_and_two_skus",
			userId:     model.UserID(1),
			want:       nil,
			cart:       twoItemsCart,
			updateCart: twoItemsCartUpdated,
		},
	}

	ctrl := minimock.NewController(t)
	productProvider := mock.NewProductProviderMock(ctrl)
	repository := mock.NewRepositoryMock(ctrl)
	lomsClient := mock.NewLOMSClientMock(ctrl)
	service := NewService(Config{ProductProviderRps: 10}, repository, productProvider, lomsClient)

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId := tt.userId

			repository.GetCartMock.Expect(ctx, userId).Return(tt.cart, tt.cartErr)
			repository.RemoveCartMock.Expect(ctx, userId)
			if tt.updateCart != nil {
				repository.UpdateCartMock.Expect(ctx, tt.updateCart)
			}

			err := service.RemoveCartItem(ctx, userId, sku)
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
