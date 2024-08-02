package stock

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"route256/loms/internal/model"
	"route256/loms/internal/usecase/stock/mock"
	"testing"
)

func TestCancelStocks(t *testing.T) {
	t.Parallel()

	ctrl := minimock.NewController(t)
	stockRepo := mock.NewRepositoryMock(ctrl)
	rStockRepo := mock.NewReservedStockRepositoryMock(ctrl)
	stockService := NewService(stockRepo, rStockRepo)

	tests := []struct {
		name                 string
		order                model.Order
		wantErr              error
		sku                  model.Sku
		stock                *model.Stock
		reservedStock        *model.ReservedStock
		reservedStockErr     error
		saveReservedStock    *model.ReservedStock
		saveReservedStockErr error
		stockErr             error
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
			reservedStockErr: ErrReservedStockNotFound,
			sku:              111,
			stock:            &model.Stock{Sku: 111, TotalCount: 2},
			wantErr:          ErrReservedStockNotFound,
		},
		{
			name: "Success",
			order: model.Order{
				ID:     1,
				UserID: 2,
				Status: model.OrderStatusNew,
				Items: []model.OrderItem{
					{Sku: 111, Count: 2},
				},
			},
			reservedStock:     &model.ReservedStock{Sku: 111, OrderID: 1, Count: 2, Status: model.ReservedStockStatusReserved},
			saveReservedStock: &model.ReservedStock{Sku: 111, OrderID: 1, Count: 2, Status: model.ReservedStockStatusCancelled},
			sku:               111,
			stock:             &model.Stock{Sku: 111, TotalCount: 2},
			wantErr:           nil,
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr == nil {
				stockRepo.SaveMock.Expect(ctx, tt.stock).Return(nil)
			}

			stockRepo.GetLockBySkuMock.Expect(ctx, tt.sku).Return(tt.stock, tt.stockErr)
			rStockRepo.GetLockedMock.Expect(ctx, tt.order.ID, tt.sku).Return(tt.reservedStock, tt.reservedStockErr)
			rStockRepo.SaveMock.Expect(ctx, tt.saveReservedStock).Return(tt.saveReservedStockErr)
			err := stockService.CancelStocks(ctx, tt.order)

			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
