package cart

import (
	"context"
	"route256/cart/internal/model"
)

func (s *Service) CheckoutCart(ctx context.Context, userId model.UserID) (model.OrderID, error) {
	ctx, span := s.beginSpan(ctx, "CheckoutCart")
	defer span.End()

	if !s.repository.HasCart(ctx, userId) {
		return 0, ErrCartNotFound
	}

	cart, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return 0, err
	}

	orderID, err := s.lomsClient.CreateOrder(ctx, cart)
	if err != nil {
		return 0, err
	}

	s.repository.RemoveCart(ctx, userId)

	return orderID, nil
}
