package noop

import (
	"context"
	"github.com/google/uuid"
	"route256/transactionalbox/pkg/transactionalbox"
	"sync"
)

type ProducerBroker struct {
	successCh chan uuid.UUID
	errorCh   chan transactionalbox.ErrSending
}

func NewProducerBroker() *ProducerBroker {
	return &ProducerBroker{
		successCh: make(chan uuid.UUID, 1),
		errorCh:   make(chan transactionalbox.ErrSending, 1),
	}
}

func (b *ProducerBroker) Send(_ context.Context, recordId uuid.UUID, _ transactionalbox.Message) error {
	b.successCh <- recordId

	return nil
}

func (b *ProducerBroker) Success(ctx context.Context) <-chan uuid.UUID {
	sync.OnceFunc(func() {
		go func() {
			defer close(b.successCh)

			for {
				select {
				case <-ctx.Done():
					return
				}
			}
		}()
	})()

	return b.successCh

}

func (b *ProducerBroker) Errors(ctx context.Context) <-chan transactionalbox.ErrSending {
	sync.OnceFunc(func() {
		b.errorCh = make(chan transactionalbox.ErrSending)

		go func() {
			defer close(b.errorCh)

			for {
				select {
				case <-ctx.Done():
					return
				}
			}
		}()
	})()

	return b.errorCh
}
