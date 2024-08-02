package cart

import (
	"cmp"
	"context"
	"route256/cart/internal/model"
	"route256/cart/internal/util/errctrl"
	"route256/cart/internal/util/rppctrl"
	"slices"
	"sync"
	"time"
)

func (s *Service) GetPublicCart(ctx context.Context, userId model.UserID) (*model.PublicCart, error) {
	ctx, span := s.beginSpan(ctx, "GetPublicCart")
	defer span.End()

	cart, err := s.repository.GetCart(ctx, userId)
	if err != nil {
		return nil, err
	}

	publicCart := &model.PublicCart{}
	cartMu := sync.Mutex{}

	errCtrl, errCtx := errctrl.NewCtrl(ctx)

	producer := func() <-chan func(context.Context) {
		ch := make(chan func(context.Context))

		go func() {
			defer close(ch)

			for _, cartItem := range cart.Items {
				wrappedFn, ok := errCtrl.Wrap(func(ctx context.Context) error {
					product, prodServErr := s.productProvider.Get(ctx, cartItem.Sku)
					if prodServErr != nil {
						return prodServErr
					}
					pubCartItem := model.PublicCartItem{Product: product, CartItem: cartItem}

					cartMu.Lock()
					defer cartMu.Unlock()

					publicCart.TotalPrice += model.Price(cartItem.Count) * product.Price
					publicCart.CartItems = append(publicCart.CartItems, pubCartItem)

					return nil
				})

				if !ok {
					break
				}

				ch <- wrappedFn
			}
		}()

		return ch
	}

	rpsCtrl := rppctrl.NewCtrl(errCtx, s.config.ProductProviderRps, time.Second)
	fns := producer()
	rpsCtrl.GoAll(fns)
	rpsCtrl.Wait()

	err = errCtrl.Wait()

	if err != nil {
		return nil, err
	}

	slices.SortFunc(publicCart.CartItems, func(a, b model.PublicCartItem) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return publicCart, nil
}
