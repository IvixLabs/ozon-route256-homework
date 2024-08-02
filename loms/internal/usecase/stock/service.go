package stock

import (
	"context"
	"route256/loms/internal/model"
	"sync"
)

type Service struct {
	stockRepository         Repository
	reservedStockRepository ReservedStockRepository
	operationMu             sync.Mutex
}

func NewService(stockRepository Repository, reservedStockRepository ReservedStockRepository) *Service {
	return &Service{
		stockRepository:         stockRepository,
		reservedStockRepository: reservedStockRepository,
	}
}

func (s *Service) doTransactionOperation(
	ctx context.Context, order model.Order,
	stockOperation func(orderItem model.OrderItem, stock *model.Stock) error,
) error {

	for _, item := range order.Items {
		stockObj, err := s.stockRepository.GetLockBySku(ctx, item.Sku)
		if err != nil {
			return err
		}

		err = stockOperation(item, stockObj)
		if err != nil {
			return err
		}

		err = s.stockRepository.Save(ctx, stockObj)
		if err != nil {
			return err
		}
	}

	return nil
}
