package events

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	edatlog "github.com/rezaAmiri123/edatV2/log"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	logger := edatlog.DefaultLogger
	handler := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.AggregateEvent])
		logger.Debug("run domainHandlers.HandleEvent")
		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

	logger.Debug("registered domain event handler transaction")
	RegisterDomainEventHandlers(subscriber, handler)
}
