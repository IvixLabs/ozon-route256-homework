package loms

import (
	"context"
	"route256/loms/internal/pb/loms/v1"
	"route256/loms/internal/usecase/order"
	"route256/loms/internal/usecase/stock"
)

var _ loms.LomsServer = (*Service)(nil)

type LomsService interface {
	OrderInfo(ctx context.Context, request *loms.OrderInfoRequest) (*loms.OrderInfoResponse, error)
	OrderCreate(ctx context.Context, request *loms.OrderCreateRequest) (*loms.OrderCreateResponse, error)
	OrderPay(context.Context, *loms.OrderPayRequest) (*loms.OrderPayResponse, error)
	OrderCancel(context.Context, *loms.OrderCancelRequest) (*loms.OrderCancelResponse, error)
	StockInfo(context.Context, *loms.StockInfoRequest) (*loms.StockInfoResponse, error)
}

type Service struct {
	loms.UnimplementedLomsServer
	impl            LomsService
	orderService    *order.Service
	stockRepository stock.Repository
	stockService    *stock.Service
}

func NewService(impl LomsService,
	orderService *order.Service,
	stockRepository stock.Repository,
	stockService *stock.Service) *Service {
	return &Service{impl: impl, orderService: orderService, stockRepository: stockRepository, stockService: stockService}
}
