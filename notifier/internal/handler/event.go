package handler

import (
	"context"
	"github.com/bufbuild/protovalidate-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/protobuf/proto"
	"route256/logger/pkg/logger"
	"route256/notifier/internal/pb/loms/v1"
	"route256/transactionalbox/pkg/transactionalbox"
	"time"
)

var msgCounterVec = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "app_notifier",
		Name:      "msg_total_counter",
		Help:      "Total amount of sql queries",
	},
	[]string{"status"},
)

var msgHistogramVec = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "app_notifier",
		Name:      "msg_duration_histogram",
		Buckets:   prometheus.DefBuckets,
	},
	[]string{"status"})

type Event struct {
}

var _ transactionalbox.ConsumerHandler = (*Event)(nil)

func (h *Event) Handle(ctx context.Context, payload []byte) error {
	ctxErr := ctx.Err()
	if ctxErr != nil {
		return ctxErr
	}

	var event *loms.Event

	defer func(startTime time.Time) {
		if event == nil {
			return
		}

		msgCounterVec.WithLabelValues(event.Status).Inc()
		msgHistogramVec.WithLabelValues(event.Status).Observe(time.Since(startTime).Seconds())
	}(time.Now())

	event = &loms.Event{}

	err := proto.Unmarshal(payload, event)
	if err != nil {
		return err
	}

	v, err := protovalidate.New()
	if err != nil {
		return err
	}

	err = v.Validate(event)
	if err != nil {
		return err
	}

	logger.Infow(ctx, "Event",
		"OrderId", event.OrderID,
		"Status", event.Status,
		"Time", event.CreatedAt.AsTime().Format(time.DateTime),
	)

	return nil
}
