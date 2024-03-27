package am

type (
	Subscription interface {
		Unsubscribe() error
	}

	MessageSubscriber interface {
		Subscription
		Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error)
	}

	messageSubscriber struct {
		subscriber MessageSubscriber
		mws        []MessageHandlerMiddleware
	}
)

var _ MessageSubscriber = (*messageSubscriber)(nil)

func NewMessageSubscriber(subscriber MessageSubscriber, mws ...MessageHandlerMiddleware) messageSubscriber {
	return messageSubscriber{
		subscriber: subscriber,
		mws:        mws,
	}
}

func (s messageSubscriber) Subscribe(topicName string, handler MessageHandler, options ...SubscriberOption) (Subscription, error) {
	return s.subscriber.Subscribe(topicName, MessageHandlerWithMiddleware(handler, s.mws...), options...)
}

func (s messageSubscriber) Unsubscribe() error {
	return s.subscriber.Unsubscribe()
}

