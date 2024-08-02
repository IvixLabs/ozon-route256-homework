package loms

import (
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/pb/loms/v1"
)

func (c *GrpcClient) CreateOrder(ctx context.Context, cart *model.Cart) (model.OrderID, error) {

	reqItems := make([]*loms.OrderItem, 0, len(cart.Items))

	for _, cartItem := range cart.Items {
		reqItems = append(reqItems,
			&loms.OrderItem{
				Sku:   uint32(cartItem.Sku),
				Count: uint32(cartItem.Count),
			})
	}

	res, err := c.client.OrderCreate(ctx, &loms.OrderCreateRequest{
		User:  int64(cart.UserId),
		Items: reqItems,
	})
	if err != nil {
		return 0, err
	}

	return model.OrderID(res.OrderID), nil
}
