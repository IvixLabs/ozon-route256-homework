package transactionalbox

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Publisher struct {
	store PublisherStore
}

type PublisherStore interface {
	AddRecord(ctx context.Context, record Record) error
}

func NewPublisher(store PublisherStore) *Publisher {
	return &Publisher{store: store}
}

func (p *Publisher) Send(ctx context.Context, msg Message) error {
	id := uuid.New()

	record := Record{
		ID:        id,
		Message:   msg,
		State:     PendingState,
		CreatedAt: time.Now(),
	}

	err := p.store.AddRecord(ctx, record)
	if err != nil {
		return err
	}

	return nil
}
