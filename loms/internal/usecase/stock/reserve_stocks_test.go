package stock

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256/loms/internal/model"
	"route256/loms/internal/usecase/stock/mock"
	"testing"
)

func TestReserveStocks(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	stockRepo := mock.NewRepositoryMock(ctrl)
	rStockRepo := mock.NewReservedStockRepositoryMock(ctrl)
	stockService := NewService(stockRepo, rStockRepo)

	errUnexpected := errors.New("unexpected error")

	tests := []struct {
		name                 string
		order                model.Order
		wantErr              error
		sku                  model.Sku
		stock                *model.Stock
		stockErr             error
		saveStock            *model.Stock
		saveStockErr         error
		wantRollback         bool
		wantRollbackError    error
		rollbackErr          error
		rollbackSaveStock    *model.Stock
		rollbackSaveStockErr error
		doRollback           bool
		isSaveReservedStock  bool
		saveReservedStock    *model.ReservedStock
		saveReservedStockErr error
	}{
		{
			name: "Get_stock_unexpected_error",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			sku:      111,
			stockErr: errUnexpected,
			wantErr:  errUnexpected,
		},
		{
			name: "Get_stock_reserve_insufficient count",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			sku:     111,
			stock:   model.NewStock(111, 1),
			wantErr: model.ErrInsufficientStockCount,
		},
		{
			name: "Get_stock_save_stock_unexpected_error",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			sku:   111,
			stock: model.NewStock(111, 2),
			saveStock: &model.Stock{
				Sku:        111,
				TotalCount: 0,
			},
			isSaveReservedStock: true,
			saveReservedStock:   &model.ReservedStock{Sku: 111, OrderID: 1, Count: 2, Status: model.ReservedStockStatusReserved},
			saveStockErr:        errUnexpected,
			wantErr:             errUnexpected,
		},
		{
			name: "Get_stock_save_stock_rollback",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			sku:   111,
			stock: model.NewStock(111, 2),
			saveStock: &model.Stock{
				Sku:        111,
				TotalCount: 0,
			},
			rollbackSaveStock: model.NewStock(111, 2),
			wantRollback:      true,
			doRollback:        true,
		},
		{
			name: "Get_stock_save_stock_rollback_error",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			sku:   111,
			stock: model.NewStock(111, 2),
			saveStock: &model.Stock{
				Sku:        111,
				TotalCount: 0,
			},
			rollbackSaveStock:    model.NewStock(111, 2),
			rollbackSaveStockErr: errUnexpected,
			wantRollbackError:    errUnexpected,
			wantRollback:         true,
			doRollback:           true,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			stockRepo.GetLockBySkuMock.Expect(ctx, tt.sku).Return(tt.stock, tt.stockErr)

			if tt.saveStock != nil {
				stockRepo.SaveMock.Expect(ctx, tt.saveStock).Return(tt.saveStockErr)
			}

			if tt.isSaveReservedStock {
				rStockRepo.SaveMock.Expect(ctx, tt.saveReservedStock).Return(tt.saveReservedStockErr)
			}

			err := stockService.ReserveStocks(ctx, tt.order)

			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
