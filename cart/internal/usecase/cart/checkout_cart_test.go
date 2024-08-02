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

func TestCheckoutCart(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	productProvider := mock.NewProductProviderMock(ctrl)
	cartRepository := mock.NewRepositoryMock(ctrl)
	lomsClient := mock.NewLOMSClientMock(ctrl)
	service := NewService(Config{ProductProviderRps: 10}, cartRepository, productProvider, lomsClient)

	ctx := context.Background()

	errUnexpected := errors.New("unexpected error")

	tests := []struct {
		name        string
		userID      model.UserID
		hasCart     bool
		wantOrderID model.OrderID
		wantErr     error
		cart        *model.Cart
		cartErr     error
		orderID     model.OrderID
		orderErr    error
	}{
		{
			name:        "Cart_not_found",
			userID:      111,
			hasCart:     false,
			wantOrderID: 0,
			wantErr:     ErrCartNotFound,
		},
		{
			name:        "Cart_found_unexpected_error",
			userID:      111,
			hasCart:     true,
			wantOrderID: 0,
			wantErr:     errUnexpected,
			cart:        &model.Cart{},
			cartErr:     errUnexpected,
		},
		{
			name:        "Cart_found_unexpected_error",
			userID:      111,
			hasCart:     true,
			wantOrderID: 0,
			wantErr:     errUnexpected,
			cart:        &model.Cart{},
			cartErr:     errUnexpected,
		},
		{
			name:        "Cart_found_create_order_unexpected_error",
			userID:      111,
			hasCart:     true,
			wantOrderID: 0,
			wantErr:     errUnexpected,
			cart:        &model.Cart{},
			orderErr:    errUnexpected,
		},
		{
			name:        "Cart_found_create_order",
			userID:      111,
			hasCart:     true,
			wantOrderID: 1,
			cart:        &model.Cart{},
			orderID:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cartRepository.HasCartMock.Expect(ctx, tt.userID).Return(tt.hasCart)
			cartRepository.GetCartMock.Expect(ctx, tt.userID).Return(tt.cart, tt.cartErr)
			lomsClient.CreateOrderMock.Expect(ctx, tt.cart).Return(tt.orderID, tt.orderErr)
			cartRepository.RemoveCartMock.Expect(ctx, tt.userID).Return()

			orderID, err := service.CheckoutCart(ctx, tt.userID)

			assert.Equal(t, tt.wantOrderID, orderID)
			assert.ErrorIs(t, tt.wantErr, err)

		})
	}

}
