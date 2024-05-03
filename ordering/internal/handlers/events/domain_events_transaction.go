package events

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/ordering/internal/constants"
)

func RegisterDomainEventHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])
		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	RegisterDomainEventHandlers(subscriber, handlers)
}
