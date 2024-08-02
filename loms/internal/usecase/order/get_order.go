package order

import (
	"context"
	"route256/loms/internal/model"
)

func (s *Service) GetOrderByID(ctx context.Context, orderId model.OrderID) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "GetOrderByID")
	defer span.End()

	order, err := s.orderRepository.GetByID(ctx, orderId)

	if err != nil {
		return nil, err
	}

	return order, nil
}
