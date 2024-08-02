package handler

import (
	"net/http"
)

func (s *Router) RemoveCart(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId, err := getUserId(r)
	if err != nil {
		return err
	}

	err = s.cartService.RemoveCart(ctx, userId)
	err = getCartNotFoundError(err)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
