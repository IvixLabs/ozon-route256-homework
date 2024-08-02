package cart

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart/mock"
	"testing"
)

func TestGetPublicCartTable(t *testing.T) {
	var (
		userId    = model.UserID(1)
		firstSku  = model.Sku(111)
		secondSku = model.Sku(222)
		thirdSku  = model.Sku(333)
		wrongSku  = model.Sku(444)

		firstProduct  = model.Product{Name: "Product 111", Price: 100}
		secondProduct = model.Product{Name: "Product 222", Price: 200}
		thirdProduct  = model.Product{Name: "Product 333", Price: 300}
	)

	createCart := func(useWrongSku bool) *model.Cart {
		cart := model.Cart{
			UserId: userId,
			Items: map[model.Sku]model.CartItem{
				firstSku:  {Sku: firstSku, Count: 1},
				thirdSku:  {Sku: thirdSku, Count: 2},
				secondSku: {Sku: secondSku, Count: 3},
			},
		}

		if useWrongSku {
			cart.Items[wrongSku] = model.CartItem{Sku: wrongSku, Count: 4}
		}

		return &cart
	}

	createPubCart := func() *model.PublicCart {
		pubCart := model.PublicCart{}

		pubCart.CartItems = append(pubCart.CartItems, model.PublicCartItem{
			Product: firstProduct, CartItem: model.CartItem{Sku: firstSku, Count: 1},
		})
		pubCart.CartItems = append(pubCart.CartItems, model.PublicCartItem{
			Product: secondProduct, CartItem: model.CartItem{Sku: secondSku, Count: 3},
		})
		pubCart.CartItems = append(pubCart.CartItems, model.PublicCartItem{
			Product: thirdProduct, CartItem: model.CartItem{Sku: thirdSku, Count: 2},
		})

		pubCart.TotalPrice = 1300

		return &pubCart
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		cartErr     error
		cart        *model.Cart
		wantErr     error
		wantPubCart *model.PublicCart
	}{
		{
			name:        "Has_not_cart",
			cartErr:     ErrCartNotFound,
			wantErr:     ErrCartNotFound,
			wantPubCart: nil,
		},
		{
			name:        "Has_cart",
			cartErr:     nil,
			wantErr:     nil,
			cart:        createCart(false),
			wantPubCart: createPubCart(),
		},
		{
			name:        "Wrong_sku",
			cartErr:     nil,
			wantErr:     ErrProductSkuNotFound,
			cart:        createCart(true),
			wantPubCart: nil,
		},
	}

	ctrl := minimock.NewController(t)
	productProvider := mock.NewProductProviderMock(ctrl)
	repository := mock.NewRepositoryMock(ctrl)
	lomsClient := mock.NewLOMSClientMock(ctrl)
	service := NewService(Config{ProductProviderRps: 10}, repository, productProvider, lomsClient)

	productProvider.
		GetMock.When(minimock.AnyContext, firstSku).Then(firstProduct, nil).
		GetMock.When(minimock.AnyContext, secondSku).Then(secondProduct, nil).
		GetMock.When(minimock.AnyContext, thirdSku).Then(thirdProduct, nil).
		GetMock.When(minimock.AnyContext, wrongSku).Then(EmptyProduct, ErrProductSkuNotFound)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)

			repository.GetCartMock.Expect(ctx, userId).Return(tt.cart, tt.cartErr)

			pubCart, err := service.GetPublicCart(ctx, userId)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, pubCart, tt.wantPubCart)
		})
	}
}
