package transactionalbox

import (
	"context"
	"github.com/google/uuid"
	"log"
	"route256/logger/pkg/logger"
	"time"
)

type ErrSending struct {
	Err error
	Key uuid.UUID
}

type ProducerBroker interface {
	Send(ctx context.Context, recordId uuid.UUID, msg Message) error
	Errors(ctx context.Context) <-chan ErrSending
	Success(ctx context.Context) <-chan uuid.UUID
}

type Producer struct {
	broker ProducerBroker
	store  ProducerStore
}

type ProducerStore interface {
	FindLockedPendingRecords(ctx context.Context, size int) ([]Record, error)
	SetDeliveredState(ctx context.Context, id uuid.UUID) error
	SetErrorState(ctx context.Context, id uuid.UUID) error
	SetPendingState(ctx context.Context, id uuid.UUID) error
}

func NewProducer(broker ProducerBroker, store ProducerStore) *Producer {
	return &Producer{broker: broker, store: store}
}

func (d *Producer) Run(ctx context.Context) {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			d.processPendingRecords(ctx)
			timer.Reset(time.Second)
		case err := <-d.broker.Errors(ctx):
			d.handleError(ctx, err)
		case suc := <-d.broker.Success(ctx):
			d.handleSuccess(ctx, suc)
		case <-ctx.Done():
			return
		}
	}
}

func (d *Producer) handleError(ctx context.Context, err ErrSending) {
	logger.Warnw(ctx, err.Err.Error())
	setErr := d.store.SetPendingState(ctx, err.Key)
	if setErr != nil {
		log.Panicln(setErr)
	}
}

func (d *Producer) handleSuccess(ctx context.Context, key uuid.UUID) {
	setErr := d.store.SetDeliveredState(ctx, key)
	if setErr != nil {
		log.Panicln(setErr)
	}
}

func (d *Producer) processPendingRecords(ctx context.Context) {
	records, err := d.store.FindLockedPendingRecords(ctx, 10)
	if err != nil {
		log.Panicln(err)
	}

	for _, rec := range records {
		sendErr := d.broker.Send(ctx, rec.ID, rec.Message)
		if sendErr != nil {
			log.Println(sendErr)
			log.Printf("RunOutboxProcessing error: %v\n", rec)
		}
	}
}
