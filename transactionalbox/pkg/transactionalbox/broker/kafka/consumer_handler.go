package kafka

import (
	"github.com/IBM/sarama"
	"route256/transactionalbox/pkg/transactionalbox"
)

type Handler struct {
	requests chan transactionalbox.ConsumerBrokerRequest
}

func NewHandler() *Handler {
	return &Handler{requests: make(chan transactionalbox.ConsumerBrokerRequest)}
}

type request struct {
	msg *sarama.ConsumerMessage
	ses sarama.ConsumerGroupSession
}

func (r *request) Payload() []byte {
	return r.msg.Value
}

func (r *request) Mark() {
	r.ses.MarkMessage(r.msg, "")

}

var _ sarama.ConsumerGroupHandler = (*Handler)(nil)

func (gh *Handler) Setup(ses sarama.ConsumerGroupSession) error {

	return nil
}

func (gh *Handler) Cleanup(ses sarama.ConsumerGroupSession) error {
	return nil
}

func (gh *Handler) ConsumeClaim(ses sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			gh.requests <- &request{msg, ses}

		case <-ses.Context().Done():
			return nil

		}
	}
}

func (h *Handler) Close() {
	close(h.requests)
}
