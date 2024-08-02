package order

import (
	"context"
	"route256/loms/internal/model"
)

func (s *Service) CancelOrder(ctx context.Context, orderId model.OrderID) error {
	ctx, span := beginSpan(ctx, "CancelOrder")
	defer span.End()

	err := s.transactionManager.Do(ctx, func(ctx context.Context) error {
		order, err := s.orderRepository.GetLockByID(ctx, orderId)
		if err != nil {
			return err
		}

		err = order.SetCanceledPaymentStatus()
		if err != nil {
			return err
		}

		err = s.stockService.CancelStocks(ctx, *order)
		if err != nil {
			return err
		}

		_, err = s.orderRepository.Save(ctx, order)
		if err != nil {
			return err
		}

		eventErr := s.sendEvent(ctx, order)
		if eventErr != nil {
			return eventErr
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
