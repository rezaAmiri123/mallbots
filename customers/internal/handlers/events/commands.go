package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(reg registry.Registry, app application.App, replyPublisher am.ReplyPublisher, amSerializer am.MessageSerializer, mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewCommandHandler(reg, amSerializer, replyPublisher, commandHandlers{
		app: app,
	}, mws...)
}

func RegisterCommandHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(customerspb.CommandChannel, handlers, am.MessageFilter{
		customerspb.AuthorizeCustomerCommand,
	}, am.GroupName("customer-commands"))
	return err
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (reply ddd.Reply, err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling command",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled command", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling command", trace.WithAttributes(
		attribute.String("Command", cmd.CommandName()),
	))

	switch cmd.CommandName() {
	case customerspb.AuthorizeCustomerCommand:
		return h.doAuthorizeCustomer(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doAuthorizeCustomer(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*customerspb.AuthorizeCustomer)

	return nil, h.app.AuthorizeCustomer(ctx, application.AuthorizeCustomer{ID: payload.GetId()})
}
