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

func TestOrderService_PayOrder(t *testing.T) {
	t.Parallel()

	unexpectedErr := errors.New("unexpected error")

	tests := []struct {
		name             string
		orderID          model.OrderID
		foundOrder       *model.Order
		foundOrderErr    error
		wantErr          error
		payStockOrder    *model.Order
		payStockRollback func() error
		payStocksErr     error
		orderSaveErr     error
	}{
		{
			name:          "Order_not_found",
			orderID:       111,
			foundOrderErr: unexpectedErr,
			payStockOrder: &model.Order{},
			wantErr:       unexpectedErr,
		},
		{
			name:          "Order_set_payed_error",
			orderID:       111,
			payStockOrder: &model.Order{},
			foundOrder:    &model.Order{Status: model.OrderStatusPayed},
			wantErr:       model.ErrWrongStatus,
		},
		{
			name:          "Order_pay_stocks_error",
			orderID:       111,
			foundOrder:    &model.Order{Status: model.OrderStatusAwaitingPayment},
			payStockOrder: &model.Order{Status: model.OrderStatusPayed},
			payStocksErr:  unexpectedErr,
			wantErr:       unexpectedErr,
		},
		{
			name:          "Order_order_save_error",
			orderID:       111,
			foundOrder:    &model.Order{Status: model.OrderStatusAwaitingPayment},
			payStockOrder: &model.Order{Status: model.OrderStatusPayed},
			payStocksErr:  nil,
			payStockRollback: func() error {
				return nil
			},
			orderSaveErr: unexpectedErr,
			wantErr:      unexpectedErr,
		},
		{
			name:          "Order_order_save_error_rollback_error",
			orderID:       111,
			foundOrder:    &model.Order{Status: model.OrderStatusAwaitingPayment},
			payStockOrder: &model.Order{Status: model.OrderStatusPayed},
			payStocksErr:  nil,
			payStockRollback: func() error {
				return unexpectedErr
			},
			orderSaveErr: unexpectedErr,
			wantErr:      unexpectedErr,
		},
		{
			name:          "Order_save_success",
			orderID:       111,
			foundOrder:    &model.Order{Status: model.OrderStatusAwaitingPayment},
			payStockOrder: &model.Order{Status: model.OrderStatusPayed},
		},
	}

	ctrl := minimock.NewController(t)
	orderRepository := mock.NewRepositoryMock(ctrl)
	stockService := mock.NewStockServiceMock(ctrl)
	txController := mock2.NewControllerMock(ctrl)
	txManager := transaction.NewManager([]transaction2.Controller{txController})
	publisher := mock.NewPublisherMock(t)
	orderService := NewService(orderRepository, publisher, stockService, txManager)
	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tx := mock2.NewTransactionMock(ctrl)

			if tt.wantErr != nil {
				tx.RollbackMock.Expect(ctx).Return(nil)
			} else {
				tx.CommitMock.Expect(ctx).Return(nil)
			}

			txController.BeginMock.Expect(ctx).Return(tx, ctx, nil)
			publisher.SendMock.Return(nil)
			orderRepository.GetLockByIDMock.Expect(ctx, tt.orderID).Return(tt.foundOrder, tt.foundOrderErr)
			stockService.PayStocksMock.Expect(ctx, *tt.payStockOrder).Return(tt.payStocksErr)
			orderRepository.SaveMock.Expect(ctx, tt.payStockOrder).Return(nil, tt.orderSaveErr)

			err := orderService.PayOrder(ctx, tt.orderID)

			assert.ErrorIs(t, tt.wantErr, err)
		})

	}

}
