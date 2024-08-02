package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_New(t *testing.T) {
	t.Parallel()

	order := NewOrder(11, []OrderItem{{Sku: 123, Count: 2}})
	expectOrder := &Order{
		UserID: 11,
		Items: []OrderItem{
			{Sku: 123, Count: 2},
		},
		Status: OrderStatusNew,
	}
	assert.Equal(t, expectOrder, order)
}

func TestOrder_SetStatus(t *testing.T) {
	t.Parallel()

	order := NewOrder(11, []OrderItem{{Sku: 123, Count: 2}})

	err := order.SetAwaitingPaymentStatus()
	assert.NoError(t, err)

	err = order.SetAwaitingPaymentStatus()
	assert.ErrorIs(t, ErrWrongStatus, err)

	err = order.SetPayedStatus()
	assert.NoError(t, err)

	err = order.SetPayedStatus()
	assert.ErrorIs(t, ErrWrongStatus, err)

	order.Status = OrderStatusAwaitingPayment

	err = order.SetCanceledPaymentStatus()
	assert.NoError(t, err)

	err = order.SetCanceledPaymentStatus()
	assert.ErrorIs(t, ErrWrongStatus, err)

	order.Status = OrderStatusNew

	err = order.SetFailedStatus()
	assert.NoError(t, err)

	err = order.SetFailedStatus()
	assert.ErrorIs(t, ErrWrongStatus, err)

}
