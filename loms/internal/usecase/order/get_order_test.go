package order

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	transaction2 "route256/common/pkg/manager/transaction"
	mock2 "route256/common/pkg/manager/transaction/mock"
	"route256/loms/internal/manager/transaction"
	"route256/loms/internal/model"
	"route256/loms/internal/usecase/order/mock"
	"testing"
)

func TestService_GetOrderByID(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	orderRepository := mock.NewRepositoryMock(ctrl)
	stockService := mock.NewStockServiceMock(ctrl)
	txController := mock2.NewControllerMock(ctrl)
	txManager := transaction.NewManager([]transaction2.Controller{txController})
	publisher := mock.NewPublisherMock(t)
	service := NewService(orderRepository, publisher, stockService, txManager)
	ctx := context.Background()

	orderRepository.GetByIDMock.Expect(ctx, 111).Return(nil, ErrOrderNotFound)
	order, err := service.GetOrderByID(ctx, 111)
	assert.Nil(t, order)
	assert.ErrorIs(t, err, ErrOrderNotFound)

	orderRepository.GetByIDMock.Expect(ctx, 222).Return(&model.Order{ID: 222}, nil)
	order, err = service.GetOrderByID(ctx, 222)
	assert.NoError(t, err)
	assert.Equal(t, order, &model.Order{ID: 222})
}
