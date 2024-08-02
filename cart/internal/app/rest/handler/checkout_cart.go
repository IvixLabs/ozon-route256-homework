package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/middleware"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
)

type checkoutCartResponse struct {
	OrderID model.OrderID `json:"orderID"`
}

func (s *Router) CheckoutCart(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userID, err := getUserId(r)
	if err != nil {
		return err
	}

	orderId, err := s.cartService.CheckoutCart(ctx, userID)
	if err != nil {
		if errors.Is(err, cart.ErrCartNotFound) {
			return fmt.Errorf("%w: %v", middleware.ErrEntityNotFound, err)
		}

		return err
	}

	bytesGetCart, err := json.Marshal(checkoutCartResponse{OrderID: orderId})
	if err != nil {
		return err
	}

	_, err = w.Write(bytesGetCart)
	if err != nil {
		return err
	}

	return nil
}
