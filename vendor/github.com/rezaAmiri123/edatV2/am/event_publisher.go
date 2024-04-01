package am

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type (
	EventPublisher interface {
		Publish(ctx context.Context, topicName string, event ddd.Event) error
	}

	eventPublisher struct {
		reg        registry.Registry
		serializer MessageSerializer
		publisher  MessagePublisher
	}
)

var _ EventPublisher = (*eventPublisher)(nil)

func NewEventPublisher(reg registry.Registry, msgPublisher MessagePublisher, serializer MessageSerializer, mws ...MessagePublisherMiddleware) EventPublisher {
	return eventPublisher{
		reg:        reg,
		serializer: serializer,
		publisher:  MessagePublisherWithMiddleware(msgPublisher, mws...),
	}
}

func (s eventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	payload, err := s.reg.Serialize(event.EventName(), event.Payload())
	if err != nil {
		return err
	}

	data, err := s.serializer.Serialize(SerializerMessageData{
		Payload:    payload,
		OccurredAt: event.OccurredAt(),
	})
	if err != nil {
		return err
	}

	return s.publisher.Publish(ctx, topicName, message{
		id:       event.ID(),
		name:     event.EventName(),
		subject:  topicName,
		data:     data,
		metadata: event.Metadata(),
		sentAt:   time.Now(),
	})
}
