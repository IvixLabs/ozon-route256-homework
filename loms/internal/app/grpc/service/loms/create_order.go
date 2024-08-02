package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/logger/pkg/logger"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
)

func (s *Service) OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error) {

	items := repackProtoToOrderItems(request.Items)
	order, err := s.orderService.CreateOrder(ctx, model.UserID(request.User), items)
	if err != nil {
		logger.Errorw(ctx, "order create error", "error", err)
		return nil, status.Errorf(codes.FailedPrecondition, "order create: %v", err)
	}

	return &loms.OrderCreateResponse{OrderID: int64(order.ID)}, nil
}

func repackProtoToOrderItems(reqItems []*loms.OrderItem) []model.OrderItem {

	items := make([]model.OrderItem, len(reqItems))
	for i, item := range reqItems {
		items[i] = model.OrderItem{Sku: model.Sku(item.Sku), Count: model.Count(item.Count)}
	}

	return items
}
