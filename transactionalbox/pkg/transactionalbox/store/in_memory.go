package store

import (
	"context"
	"github.com/google/uuid"
	"route256/transactionalbox/pkg/transactionalbox"
	"sync"
)

var storage map[uuid.UUID]transactionalbox.Record = make(map[uuid.UUID]transactionalbox.Record)

type InMemory struct {
	storage map[uuid.UUID]transactionalbox.Record
	mu      sync.Mutex
}

var _ transactionalbox.ProducerStore = (*InMemory)(nil)

func NewInMemory() *InMemory {
	return &InMemory{storage: storage}
}

func (s *InMemory) AddRecord(_ context.Context, record transactionalbox.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[record.ID] = record

	return nil
}

func (s *InMemory) FindLockedPendingRecords(_ context.Context, size int) ([]transactionalbox.Record, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]transactionalbox.Record, 0, size)

	for key, rec := range s.storage {
		if rec.State == transactionalbox.PendingState {
			rec.State = transactionalbox.ProcessingState

			s.storage[key] = rec

			res = append(res, rec)
			if len(res) == size {
				break
			}
		}
	}

	return res, nil
}

func (s *InMemory) SetDeliveredState(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.storage[id]

	if !ok {
		return nil
	}

	rec.State = transactionalbox.DeliveredState

	s.storage[id] = rec

	return nil
}

func (s *InMemory) SetErrorState(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.storage[id]

	if !ok {
		return nil
	}

	rec.State = transactionalbox.ErrorState
	s.storage[id] = rec

	return nil
}

func (s *InMemory) SetPendingState(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.storage[id]

	if !ok {
		return nil
	}

	rec.State = transactionalbox.PendingState
	s.storage[id] = rec

	return nil
}
