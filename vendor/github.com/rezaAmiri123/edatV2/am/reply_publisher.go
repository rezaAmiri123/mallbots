package am

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type (
	ReplyPublisher interface {
		Publish(ctx context.Context, topicName string, reply ddd.Reply) error
	}

	replyPublisher struct {
		reg        registry.Registry
		serializer MessageSerializer
		publisher  MessagePublisher
	}
)

var _ ReplyPublisher = (*replyPublisher)(nil)

func NewReplyPublisher(
	reg registry.Registry,
	serializer MessageSerializer,
	publisher MessagePublisher,
	mws ...MessagePublisherMiddleware,
) ReplyPublisher {
	return replyPublisher{
		reg:        reg,
		serializer: serializer,
		publisher:  MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (s replyPublisher) Publish(ctx context.Context, topicName string, reply ddd.Reply) error {
	var err error
	var payload []byte

	if reply.ReplyName() != SuccessReply && reply.ReplyName() != FailureReply {
		payload, err = s.reg.Serialize(reply.ReplyName(), reply.Payload())
		if err != nil {
			return err
		}
	}

	data, err := s.serializer.Serialize(SerializerMessageData{
		Payload:    payload,
		OccurredAt: reply.OccurredAt(),
	})
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       reply.ID(),
		name:     reply.ReplyName(),
		subject:  topicName,
		data:     data,
		metadata: reply.Metadata(),
		sentAt:   time.Now(),
	})
}
