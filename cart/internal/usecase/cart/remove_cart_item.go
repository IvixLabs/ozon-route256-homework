package cart

import (
	"context"
	"route256/cart/internal/model"
)

func (s *Service) RemoveCartItem(ctx context.Context, userId model.UserID, sku model.Sku) error {
	ctx, span := s.beginSpan(ctx, "RemoveCartItem")
	defer span.End()

	cart, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return err
	}

	if _, ok := cart.Items[sku]; !ok {
		return ErrCartItemNotFound
	}

	if len(cart.Items) == 1 {
		s.repository.RemoveCart(ctx, userId)
	} else {
		delete(cart.Items, sku)
		s.repository.UpdateCart(ctx, cart)
	}

	return nil
}
