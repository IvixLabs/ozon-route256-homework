package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
)

var (
	statusMap = map[model.OrderStatus]string{
		model.OrderStatusNew:             "new",
		model.OrderStatusAwaitingPayment: "awaiting payment",
		model.OrderStatusFailed:          "failed",
		model.OrderStatusPayed:           "payed",
		model.OrderStatusCancelled:       "cancelled",
	}
)

func (s *Service) OrderInfo(ctx context.Context, request *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error) {

	order, err := s.orderService.GetOrderByID(ctx, model.OrderID(request.OrderID))

	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "order info: %v", err)
	}

	return repackOrderToProto(order), err
}

func repackOrderToProto(order *model.Order) *loms.OrderInfoResponse {

	items := make([]*loms.OrderItem, len(order.Items))

	for i, item := range order.Items {
		items[i] = &loms.OrderItem{Sku: uint32(item.Sku), Count: uint32(item.Count)}
	}

	res := &loms.OrderInfoResponse{
		Status: statusMap[order.Status],
		Items:  items,
	}

	return res
}
