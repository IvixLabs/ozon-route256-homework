package order

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	transaction2 "route256/common/pkg/manager/transaction"
	mock2 "route256/common/pkg/manager/transaction/mock"
	"route256/loms/internal/manager/transaction"
	"route256/loms/internal/model"
	"route256/loms/internal/usecase/order/mock"
	"testing"
)

func TestOrderService_CreateOrder(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	orderRepository := mock.NewRepositoryMock(ctrl)
	stockService := mock.NewStockServiceMock(ctrl)
	txController := mock2.NewControllerMock(ctrl)

	txManager := transaction.NewManager([]transaction2.Controller{txController})
	publisher := mock.NewPublisherMock(t)
	orderService := NewService(orderRepository, publisher, stockService, txManager)

	unexpectedErr := errors.New("unexpected error")
	ctx := context.Background()

	tests := []struct {
		name       string
		userId     model.UserID
		orderItems []model.OrderItem
		wantOrder  *model.Order
		wantErr    error
		prepare    func()
	}{
		{
			name:       "Save_order_unexpected_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantErr:    unexpectedErr,
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)
				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				orderRepository.SaveMock.Expect(minimock.AnyContext, saveOrder).Return(nil, unexpectedErr)
			},
		},
		{
			name:       "Reserve_stocks_unexpected_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				publisher.SendMock.Return(nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				retOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				orderRepository.
					SaveMock.When(ctx, saveOrder).Then(retOrder, nil)
				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(unexpectedErr)
			},
			wantOrder: nil,
			wantErr:   unexpectedErr,
		},
		{
			name:       "Reserve_stocks_insufficient_stock_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantErr:    model.ErrWrongStatus,
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				publisher.SendMock.Return(nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				retOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusFailed}
				orderRepository.SaveMock.Expect(ctx, saveOrder).Return(retOrder, nil)

				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(model.ErrInsufficientStockCount)

			},
		},
		{
			name:       "Reserve_stocks_unexpected_error_set_failed_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantErr:    unexpectedErr,
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				publisher.SendMock.Return(nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				retOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				awaitingPaymentOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusAwaitingPayment}
				orderRepository.
					SaveMock.When(ctx, saveOrder).Then(retOrder, nil).
					SaveMock.When(ctx, awaitingPaymentOrder).Then(nil, unexpectedErr)
				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(nil)

			},
		},
		{
			name:       "Reserve_stocks_insufficient_stock_error_save_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantErr:    unexpectedErr,
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				publisher.SendMock.Return(nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				newOrder := &model.Order{
					UserID: 111,
					Items:  []model.OrderItem{{Sku: 222, Count: 2}},
					Status: model.OrderStatusNew,
				}
				retOrder := &model.Order{
					ID:     1,
					UserID: 111,
					Items:  []model.OrderItem{{Sku: 222, Count: 2}},
					Status: model.OrderStatusNew,
				}

				failedOrder := &model.Order{
					ID:     1,
					UserID: 111,
					Items:  []model.OrderItem{{Sku: 222, Count: 2}},
					Status: model.OrderStatusFailed,
				}
				orderRepository.
					SaveMock.When(ctx, newOrder).Then(retOrder, nil).
					SaveMock.When(ctx, failedOrder).Then(nil, unexpectedErr)

				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(model.ErrInsufficientStockCount)

			},
		},
		{
			name:       "Reserve_stocks_unexpected_error_set_awaiting_payment_error",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantErr:    model.ErrWrongStatus,
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.RollbackMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				publisher.SendMock.Return(nil)
				tx.CommitMock.Expect(ctx).Return(nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				retOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusAwaitingPayment}
				orderRepository.
					SaveMock.Expect(ctx, saveOrder).Return(retOrder, nil)

				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(nil)
			},
		},
		{
			name:       "Updated_order_save",
			userId:     111,
			orderItems: []model.OrderItem{{Sku: 222, Count: 2}},
			wantOrder:  &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusAwaitingPayment},
			prepare: func() {
				tx := mock2.NewTransactionMock(ctrl)
				tx.CommitMock.Expect(ctx).Return(nil)
				txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)

				orderRepository = mock.NewRepositoryMock(ctrl)
				stockService = mock.NewStockServiceMock(ctrl)
				orderService = NewService(orderRepository, publisher, stockService, txManager)

				saveOrder := &model.Order{UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}
				retOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusNew}

				updatedOrder := &model.Order{ID: 1, UserID: 111, Items: []model.OrderItem{{Sku: 222, Count: 2}}, Status: model.OrderStatusAwaitingPayment}
				orderRepository.
					SaveMock.When(ctx, saveOrder).Then(retOrder, nil).
					SaveMock.When(ctx, updatedOrder).Then(updatedOrder, nil)

				stockService.ReserveStocksMock.Expect(ctx, *retOrder).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			order, err := orderService.CreateOrder(ctx, tt.userId, tt.orderItems)
			assert.Equal(t, tt.wantOrder, order)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
