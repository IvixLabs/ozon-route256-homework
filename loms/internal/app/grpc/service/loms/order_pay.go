package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
)

func (s *Service) OrderPay(ctx context.Context, req *loms.OrderPayRequest) (*loms.OrderPayResponse, error) {

	err := s.orderService.PayOrder(ctx, model.OrderID(req.OrderID))

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "order pay: %v", err)
	}

	return &loms.OrderPayResponse{}, nil
}
