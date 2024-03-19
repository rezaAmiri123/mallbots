package am

type (
	MessageStream interface {
		MessagePublisher
		MessageSubscriber
	}

	MessageStreamMiddleware = func(next MessageStream)MessageStream
)
