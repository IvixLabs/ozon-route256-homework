package model

import "errors"

type ReservedStockStatus string

var (
	ReservedStockStatusReserved  = ReservedStockStatus("reserved")
	ReservedStockStatusCancelled = ReservedStockStatus("cancelled")
	ReservedStockStatusPaid      = ReservedStockStatus("paid")
)

var (
	ErrReservedStockWrongStatus = errors.New("reserved stock has wrong status")
)

type ReservedStock struct {
	OrderID OrderID
	Sku     Sku
	Count   Count
	Status  ReservedStockStatus
}

func NewReservedStock(orderID OrderID, sku Sku, count Count) *ReservedStock {
	return &ReservedStock{OrderID: orderID, Sku: sku, Count: count, Status: ReservedStockStatusReserved}
}

func (rs *ReservedStock) SetPaidStatus() error {
	if rs.Status != ReservedStockStatusReserved {
		return ErrReservedStockWrongStatus
	}
	rs.Status = ReservedStockStatusPaid

	return nil
}

func (rs *ReservedStock) SetCancelledStatus() error {
	if rs.Status != ReservedStockStatusReserved {
		return ErrReservedStockWrongStatus
	}
	rs.Status = ReservedStockStatusCancelled

	return nil
}
