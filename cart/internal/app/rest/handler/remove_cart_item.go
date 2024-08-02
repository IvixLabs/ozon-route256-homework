package handler

import (
	"net/http"
)

func (s *Router) RemoveCartItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId, err := getUserId(r)
	if err != nil {
		return err
	}

	sku, err := getSku(r)
	if err != nil {
		return err
	}

	err = s.cartService.RemoveCartItem(ctx, userId, sku)
	errCart := getCartNotFoundError(err)
	if errCart != nil {
		return errCart
	}

	errCartItem := getCartItemNotFoundError(err)
	if errCartItem != nil {
		return errCartItem
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
