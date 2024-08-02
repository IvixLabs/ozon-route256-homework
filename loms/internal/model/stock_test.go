package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStock_New(t *testing.T) {
	t.Parallel()

	stock := NewStock(111, 5)
	assert.Equal(t, Sku(111), stock.Sku)
	assert.Equal(t, Count(5), stock.TotalCount)
}

func TestStock_ReserveCount(t *testing.T) {
	t.Parallel()

	stock := NewStock(111, 5)
	rStock, err := stock.Reserve(OrderID(222), Count(1))
	assert.NoError(t, err)
	assert.Equal(t, Count(4), stock.TotalCount)
	assert.Equal(t, &ReservedStock{OrderID: OrderID(222), Sku: Sku(111), Count: Count(1), Status: ReservedStockStatusReserved}, rStock)

	_, err = stock.Reserve(222, 5)
	assert.ErrorIs(t, err, ErrInsufficientStockCount)
}

func TestStock_CancelCount(t *testing.T) {
	t.Parallel()

	stock := NewStock(111, 5)

	wrongRStock := &ReservedStock{Status: ReservedStockStatusCancelled}
	err := stock.Cancel(wrongRStock)
	assert.ErrorIs(t, err, ErrReservedStockWrongStatus)

	rStock, err := stock.Reserve(222, 3)
	err = stock.Cancel(rStock)
	assert.NoError(t, err)
	assert.Equal(t, Count(5), stock.TotalCount)
	assert.Equal(t, &ReservedStock{OrderID: OrderID(222), Sku: Sku(111), Count: Count(3), Status: ReservedStockStatusCancelled}, rStock)
}

func TestStock_PayCount(t *testing.T) {
	t.Parallel()

	stock := NewStock(111, 5)

	wrongRStock := &ReservedStock{Status: ReservedStockStatusCancelled}
	err := stock.Pay(wrongRStock)
	assert.ErrorIs(t, err, ErrReservedStockWrongStatus)

	rStock, err := stock.Reserve(222, 3)
	err = stock.Pay(rStock)
	assert.NoError(t, err)
	assert.Equal(t, Count(2), stock.TotalCount)
	assert.Equal(t, &ReservedStock{OrderID: OrderID(222), Sku: Sku(111), Count: Count(3), Status: ReservedStockStatusPaid}, rStock)
}

func TestStock_Clone(t *testing.T) {
	t.Parallel()

	stock := NewStock(111, 3)

	clonedStock := stock.Clone()

	stock.Sku = 222
	stock.TotalCount = 5

	expectedStock := NewStock(111, 3)

	assert.Equal(t, expectedStock, &clonedStock)

}
