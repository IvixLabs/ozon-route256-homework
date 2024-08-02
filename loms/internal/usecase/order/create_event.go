package order

import (
	"context"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"route256/loms/internal/model"
	"route256/loms/internal/pb/loms/v1"
	"route256/transactionalbox/pkg/transactionalbox"
	"time"
)

const EventTopic = "loms.order-events"

func (s *Service) sendEvent(ctx context.Context, order *model.Order) error {

	event := &loms.Event{
		OrderID:   int64(order.ID),
		Status:    string(order.Status),
		CreatedAt: timestamppb.New(time.Now()),
	}

	body, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, uint64(order.ID))

	msg := transactionalbox.Message{
		Key:   key,
		Topic: EventTopic,
		Body:  body,
	}

	err = s.publisher.Send(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
