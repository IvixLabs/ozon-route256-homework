package model

import "errors"

type OrderStatus string

var (
	OrderStatusNew             = OrderStatus("new")
	OrderStatusAwaitingPayment = OrderStatus("awaiting_payment")
	OrderStatusFailed          = OrderStatus("failed")
	OrderStatusPayed           = OrderStatus("payed")
	OrderStatusCancelled       = OrderStatus("cancelled")
)

var (
	ErrWrongStatus = errors.New("wrong status")
)

type OrderItem struct {
	Sku   Sku
	Count Count
}

type Order struct {
	ID     OrderID
	UserID UserID
	Status OrderStatus
	Items  []OrderItem
}

func NewOrder(userID UserID, items []OrderItem) *Order {
	return &Order{
		UserID: userID,
		Status: OrderStatusNew,
		Items:  items,
	}
}

func (o *Order) SetPayedStatus() error {
	if o.Status != OrderStatusAwaitingPayment {
		return ErrWrongStatus
	}

	o.Status = OrderStatusPayed

	return nil
}

func (o *Order) SetAwaitingPaymentStatus() error {
	if o.Status != OrderStatusNew {
		return ErrWrongStatus
	}

	o.Status = OrderStatusAwaitingPayment

	return nil
}

func (o *Order) SetCanceledPaymentStatus() error {
	if o.Status != OrderStatusAwaitingPayment {
		return ErrWrongStatus
	}

	o.Status = OrderStatusCancelled

	return nil
}

func (o *Order) SetFailedStatus() error {
	if o.Status != OrderStatusNew {
		return ErrWrongStatus
	}

	o.Status = OrderStatusFailed

	return nil
}
