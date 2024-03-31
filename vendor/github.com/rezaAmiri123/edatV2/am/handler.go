package am

import "context"

type (
	MessageHandler interface {
		HandleMessage(ctx context.Context, msg IncomingMessage) error
	}

	MessageHandlerFunc func(ctx context.Context, msg IncomingMessage) error

	MessageHandlerMiddleware = func(next MessageHandler) MessageHandler
)

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	return f(ctx, msg)
}
