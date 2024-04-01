package am

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type (
	CommandPublisher interface {
		Publish(ctx context.Context, topicName string, cmd ddd.Command) error
	}

	commandPublisher struct {
		reg        registry.Registry
		serializer MessageSerializer
		publisher  MessagePublisher
	}
)

var _ CommandPublisher = (*commandPublisher)(nil)

func NewCommandPublisher(
	reg registry.Registry,
	serializer MessageSerializer,
	publisher MessagePublisher,
	mws ...MessagePublisherMiddleware,
) CommandPublisher {
	return commandPublisher{
		reg:        reg,
		serializer: serializer,
		publisher:  MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (c commandPublisher) Publish(ctx context.Context, topicName string, command ddd.Command) error {
	payload, err := c.reg.Serialize(command.CommandName(), command.Payload())
	if err != nil {
		return err
	}

	data, err := c.serializer.Serialize(SerializerMessageData{
		Payload:    payload,
		OccurredAt: command.OccurredAt(),
	})
	if err != nil {
		return err
	}

	return c.publisher.Publish(ctx, topicName, message{
		id:       command.ID(),
		name:     command.CommandName(),
		subject:  topicName,
		data:     data,
		metadata: command.Metadata(),
		sentAt:   time.Now(),
	})
}
