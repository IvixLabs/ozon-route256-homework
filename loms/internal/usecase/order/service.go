package order

import (
	"context"
	"route256/loms/internal/manager/transaction"
	"route256/loms/internal/model"
	"route256/transactionalbox/pkg/transactionalbox"
)

type Service struct {
	orderRepository    Repository
	publisher          Publisher
	stockService       StockService
	transactionManager *transaction.Manager
}

type StockService interface {
	ReserveStocks(ctx context.Context, order model.Order) error
	PayStocks(ctx context.Context, order model.Order) error
	CancelStocks(ctx context.Context, order model.Order) error
}

type Publisher interface {
	Send(ctx context.Context, msg transactionalbox.Message) error
}

func NewService(orderRepository Repository,
	publisher Publisher,
	stockService StockService,
	transactionManager *transaction.Manager) *Service {
	return &Service{orderRepository: orderRepository,
		publisher:          publisher,
		stockService:       stockService,
		transactionManager: transactionManager}
}
