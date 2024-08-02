package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/cart/internal/model"
	"route256/cart/internal/pb/loms/v1"
	"route256/cart/internal/usecase/cart"
)

func (c *GrpcClient) GetStockCount(ctx context.Context, sku model.Sku) (model.Count, error) {

	res, err := c.client.StockInfo(ctx, &loms.StockInfoRequest{Sku: uint32(sku)})

	if err != nil {
		stat, ok := status.FromError(err)
		if !ok {
			return 0, err
		}

		if stat.Code() == codes.FailedPrecondition {
			return 0, cart.ErrStockNotFound
		}

		return 0, err
	}

	return model.Count(res.Count), nil
}
