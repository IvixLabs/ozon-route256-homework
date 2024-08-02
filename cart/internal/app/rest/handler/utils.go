package handler

import (
	"errors"
	"fmt"
	"net/http"
	"route256/cart/internal/middleware"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
	"strconv"
)

func getUserId(r *http.Request) (model.UserID, error) {
	strUserId := r.PathValue("id")
	userId, err := strconv.ParseInt(strUserId, 10, 64)
	if err != nil || userId < 1 {
		return 0, fmt.Errorf("%w: wrong id", middleware.ErrValidation)
	}

	return model.UserID(userId), nil
}

func getSku(r *http.Request) (model.Sku, error) {
	strSku := r.PathValue("sku")
	sku, err := strconv.ParseInt(strSku, 10, 64)
	if err != nil || sku < 1 {
		return 0, fmt.Errorf("%w: wrong sku", middleware.ErrValidation)
	}

	return model.Sku(sku), nil
}

func getCartNotFoundError(err error) error {
	if err != nil {
		if errors.Is(err, cart.ErrCartNotFound) {
			return fmt.Errorf("%w: cart is not found", middleware.ErrEntityNotFound)
		}
	}

	return nil
}

func getCartItemNotFoundError(err error) error {
	if err != nil {
		if errors.Is(err, cart.ErrCartItemNotFound) {
			return fmt.Errorf("%w: cart item is not found", middleware.ErrEntityNotFound)
		}
	}

	return nil
}
