package order

import (
	"context"
	"route256/loms/internal/model"
)

func (s *Service) PayOrder(ctx context.Context, orderId model.OrderID) error {
	ctx, span := beginSpan(ctx, "PayOrder")
	defer span.End()

	err := s.transactionManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.GetLockByID(ctx, orderId)
		if err != nil {
			return err
		}

		err = order.SetPayedStatus()
		if err != nil {
			return err
		}

		err = s.stockService.PayStocks(ctx, *order)
		if err != nil {
			return err
		}

		_, err = s.orderRepository.Save(ctx, order)

		eventErr := s.sendEvent(ctx, order)
		if eventErr != nil {
			return eventErr
		}

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
