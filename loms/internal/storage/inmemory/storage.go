package inmemory

import (
	"route256/loms/internal/model"
	"sync"
)

type Storage struct {
	Orders   map[model.OrderID]model.Order
	muOrders sync.Mutex
	Stocks   map[model.Sku]model.Stock
	muStocks sync.Mutex

	ReservedStocks   map[string]model.ReservedStock
	muReservedStocks sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		Orders:         make(map[model.OrderID]model.Order),
		Stocks:         make(map[model.Sku]model.Stock),
		ReservedStocks: make(map[string]model.ReservedStock),
	}
}

func (s *Storage) LockOrders() {
	s.muOrders.Lock()
}

func (s *Storage) UnlockOrders() {
	s.muOrders.Unlock()
}

func (s *Storage) LockStocks() {
	s.muStocks.Lock()
}

func (s *Storage) UnlockStocks() {
	s.muStocks.Unlock()
}

func (s *Storage) LockReservedStocks() {
	s.muReservedStocks.Lock()
}

func (s *Storage) UnlockReservedStocks() {
	s.muReservedStocks.Unlock()
}
