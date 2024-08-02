package cart

import (
	"context"
	"route256/cart/internal/model"
)

func (s *Service) RemoveCart(ctx context.Context, userId model.UserID) error {
	ctx, span := s.beginSpan(ctx, "RemoveCart")
	defer span.End()

	if !s.repository.HasCart(ctx, userId) {
		return ErrCartNotFound
	}

	s.repository.RemoveCart(ctx, userId)

	return nil
}
