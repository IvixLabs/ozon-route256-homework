package loms

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/logger/pkg/logger"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
	"route256/loms/internal/usecase/order"
	"route256/loms/internal/usecase/stock"
)

func (s *Service) OrderCancel(ctx context.Context, req *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error) {

	err := s.orderService.CancelOrder(ctx, model.OrderID(req.OrderID))

	if err != nil {
		if errors.Is(err, stock.ErrReservedStockNotFound) || errors.Is(err, order.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "order cancel: %v", err)
		}

		logger.Errorw(ctx, "order cancel error", "error", err)
		return nil, status.Errorf(codes.Internal, "order cancel: %v", err)
	}

	return &loms.OrderCancelResponse{}, nil
}
