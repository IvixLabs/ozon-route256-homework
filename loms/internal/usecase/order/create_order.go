package order

import (
	"context"
	"errors"
	"route256/loms/internal/model"
)

func (s *Service) CreateOrder(ctx context.Context, userId model.UserID, orderItems []model.OrderItem) (*model.Order, error) {
	ctx, span := beginSpan(ctx, "CreateOrder")
	defer span.End()

	var order *model.Order
	var err error

	err = s.transactionManager.Do(ctx, func(ctx context.Context) error {
		order, err = s.orderRepository.Save(ctx, model.NewOrder(userId, orderItems))
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
		return nil, err
	}

	err = s.transactionManager.Do(ctx, func(ctx context.Context) error {
		innerErr := s.stockService.ReserveStocks(ctx, *order)
		if innerErr != nil {
			return innerErr
		}

		innerErr = order.SetAwaitingPaymentStatus()
		if innerErr != nil {
			return innerErr
		}

		order, innerErr = s.orderRepository.Save(ctx, order)
		if innerErr != nil {
			return innerErr
		}

		eventErr := s.sendEvent(ctx, order)
		if eventErr != nil {
			return eventErr
		}

		return nil

	})

	if err != nil {
		if errors.Is(err, model.ErrInsufficientStockCount) {
			statusSrr := order.SetFailedStatus()
			if statusSrr != nil {
				return nil, statusSrr
			}
			_, saveErr := s.orderRepository.Save(ctx, order)
			if saveErr != nil {
				return nil, saveErr
			}
		}

		return nil, err
	}

	return order, nil
}
