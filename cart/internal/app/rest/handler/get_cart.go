package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
)

type getCartItemResponse struct {
	SkuId model.Sku   `json:"sku_id"`
	Name  string      `json:"name"`
	Count model.Count `json:"count"`
	Price model.Price `json:"price"`
}

type getCartResponse struct {
	Items      []getCartItemResponse `json:"items"`
	TotalPrice model.Price           `json:"total_price"`
}

func (s *Router) GetCart(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId, err := getUserId(r)
	if err != nil {
		return err
	}

	pubCart, err := s.cartService.GetPublicCart(ctx, userId)
	if err != nil {
		if errors.Is(err, cart.ErrCartNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return nil
		} else {
			return err
		}
	}

	getCart := getCartResponse{
		Items:      make([]getCartItemResponse, len(pubCart.CartItems)),
		TotalPrice: pubCart.TotalPrice,
	}

	for i, cartItem := range pubCart.CartItems {
		getCart.Items[i] = getCartItemResponse{
			SkuId: cartItem.Sku,
			Name:  cartItem.Name,
			Count: cartItem.Count,
			Price: cartItem.Price,
		}
	}

	bytesGetCart, err := json.Marshal(getCart)
	if err != nil {
		return err
	}

	_, err = w.Write(bytesGetCart)
	if err != nil {
		return err
	}

	return nil
}
