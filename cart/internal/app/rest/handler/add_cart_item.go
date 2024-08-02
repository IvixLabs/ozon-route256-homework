package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"route256/cart/internal/middleware"
	"route256/cart/internal/model"
	"route256/cart/internal/usecase/cart"
)

type addCartItemRequest struct {
	Count model.Count `json:"count" validate:"gt=0"`
}

func (s *Router) AddCartItem(_ http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	userId, err := getUserId(r)
	if err != nil {
		return err
	}

	sku, err := getSku(r)
	if err != nil {
		return err
	}

	bytesReq, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	req := &addCartItemRequest{}
	err = json.Unmarshal(bytesReq, req)
	if err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	validErr := validate.Struct(req)
	if validErr != nil {
		return fmt.Errorf("%w "+validErr.Error(), middleware.ErrValidation)
	}

	err = s.cartService.AddCartItem(ctx, userId, sku, req.Count)
	if err != nil {
		switch {
		case errors.Is(err, cart.ErrProductSkuNotFound):
			return fmt.Errorf("%w: %v", middleware.ErrWrongArgument, err)
		case errors.Is(err, cart.ErrInsufficientStocks):
			return fmt.Errorf("%w: %v", middleware.ErrWrongArgument, err)
		case errors.Is(err, cart.ErrStockNotFound):
			return fmt.Errorf("%w: %v", middleware.ErrWrongArgument, err)
		default:
			return err
		}
	}

	return nil
}
