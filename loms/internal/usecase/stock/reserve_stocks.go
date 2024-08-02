package stock

import (
	"context"
	"route256/loms/internal/model"
)

func (s *Service) ReserveStocks(ctx context.Context, order model.Order) error {
	ctx, span := beginSpan(ctx, "ReserveStocks")
	defer span.End()

	err := s.doTransactionOperation(ctx, order, func(orderItem model.OrderItem, stock *model.Stock) error {

		rStock, err := stock.Reserve(order.ID, orderItem.Count)
		if err != nil {
			return err
		}

		err = s.reservedStockRepository.Save(ctx, rStock)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
