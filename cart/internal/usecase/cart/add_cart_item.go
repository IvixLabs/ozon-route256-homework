package cart

import (
	"context"
	"errors"
	"route256/cart/internal/model"
)

var (
	ErrWrongCartItemCount = errors.New("wrong cart item count")
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (s *Service) AddCartItem(ctx context.Context, userId model.UserID, sku model.Sku, count model.Count) error {
	ctx, span := s.beginSpan(ctx, "AddCartItem")
	defer span.End()

	if count < 1 {
		return ErrWrongCartItemCount
	}

	cart, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrCartNotFound) {
			cart = model.NewCart(userId)
		} else {
			return err
		}
	}

	var cartItem model.CartItem
	var ok bool
	if cartItem, ok = cart.Items[sku]; ok {
		cartItem.Count += count
	} else {
		_, err = s.productProvider.Get(ctx, sku)
		if err != nil {
			return err
		}

		cartItem = model.CartItem{
			Sku:   sku,
			Count: count,
		}
	}

	stockCount, err := s.lomsClient.GetStockCount(ctx, sku)
	if err != nil {
		return err
	}

	if stockCount < cartItem.Count {
		return ErrInsufficientStocks
	}

	cart.AddItem(cartItem)

	s.repository.UpdateCart(ctx, cart)

	return nil
}
