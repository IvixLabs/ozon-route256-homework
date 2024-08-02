package model

import "errors"

var (
	ErrInsufficientStockCount = errors.New("insufficient stock count")
)

type Stock struct {
	Sku        Sku
	TotalCount Count
}

func NewStock(sku Sku, totalCount Count) *Stock {
	return &Stock{
		Sku:        sku,
		TotalCount: totalCount,
	}
}

func (s *Stock) Reserve(orderID OrderID, count Count) (*ReservedStock, error) {
	if s.TotalCount < count {
		return nil, ErrInsufficientStockCount
	}

	s.TotalCount -= count

	rStock := NewReservedStock(orderID, s.Sku, count)

	return rStock, nil
}

func (s *Stock) Cancel(reservedStock *ReservedStock) error {

	err := reservedStock.SetCancelledStatus()
	if err != nil {
		return err
	}

	s.TotalCount += reservedStock.Count
	reservedStock.Status = ReservedStockStatusCancelled

	return nil
}

func (s *Stock) Pay(reservedStock *ReservedStock) error {

	err := reservedStock.SetPaidStatus()
	if err != nil {
		return err
	}

	reservedStock.Status = ReservedStockStatusPaid

	return nil
}

func (s *Stock) Clone() Stock {
	newStock := *s

	return newStock
}
