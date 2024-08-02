package cart

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart/mock"
	"testing"
)

func TestRemoveCartTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userId  model.UserID
		want    error
		hasCart bool
	}{
		{
			name:    "Has_cart",
			userId:  model.UserID(1),
			hasCart: true,
			want:    nil,
		},
		{
			name:    "Has_not_cart",
			userId:  model.UserID(2),
			hasCart: false,
			want:    ErrCartNotFound,
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

			repository.HasCartMock.Expect(ctx, userId).Return(tt.hasCart)
			repository.RemoveCartMock.Expect(ctx, userId)

			err := service.RemoveCart(ctx, userId)
			assert.ErrorIs(t, err, tt.want)
		})
	}
}
