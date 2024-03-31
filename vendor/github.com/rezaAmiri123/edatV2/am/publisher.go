package am

import "context"

type (
	// TODO add options like partition for publisher
	MessagePublisher interface {
		Publish(ctx context.Context, topicName string, msg Message) error
	}

	MessagePublisherFunc func(ctx context.Context, topicName string, msg Message) error

	MessagePublisherMiddleware = func(next MessagePublisher) MessagePublisher

	messagePublisher struct {
		publisher MessagePublisher
	}
)

var _ MessagePublisher = (*messagePublisher)(nil)

func NewMessagePublisher(publisher MessagePublisher, mws ...MessagePublisherMiddleware) MessagePublisher {
	return messagePublisher{
		publisher: MessagePublisherWithMiddleware(publisher, mws...),
	}
}

func (p messagePublisher) Publish(ctx context.Context, topicName string, msg Message) error {
	return p.publisher.Publish(ctx, topicName, msg)
}

func (f MessagePublisherFunc) Publish(ctx context.Context, topicName string, msg Message) error {
	return f(ctx, topicName, msg)
}
