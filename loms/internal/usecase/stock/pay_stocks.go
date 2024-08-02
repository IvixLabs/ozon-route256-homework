package stock

import (
	"context"
	"route256/loms/internal/model"
)

func (s *Service) PayStocks(ctx context.Context, order model.Order) error {
	ctx, span := beginSpan(ctx, "PayStocks")
	defer span.End()

	err := s.doTransactionOperation(ctx, order, func(_ model.OrderItem, stock *model.Stock) error {
		rStock, err := s.reservedStockRepository.GetLocked(ctx, order.ID, stock.Sku)
		if err != nil {
			return err
		}

		err = stock.Pay(rStock)
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
