package transactionalbox

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Key   []byte
	Body  []byte
	Topic string
}

type Record struct {
	ID        uuid.UUID
	Message   Message
	State     RecordState
	CreatedAt time.Time
}

type RecordState int

const (
	PendingState RecordState = iota
	ProcessingState
	DeliveredState
	ErrorState
)

type Store interface {
	PublisherStore
	ProducerStore
}
