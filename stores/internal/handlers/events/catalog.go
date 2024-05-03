package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/stores/internal/constants"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type catalagHandlers[T ddd.Event] struct {
	catalog domain.CatalogRepository
}

var _ ddd.EventHandler[ddd.Event] = (*catalagHandlers[ddd.Event])(nil)

func NewCatalagHandlers(catalog domain.CatalogRepository) ddd.EventHandler[ddd.Event] {
	return catalagHandlers[ddd.Event]{
		catalog: catalog,
	}
}

func RegisterCatalagHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.ProductAddedEvent,
	)
}

func RegisterCatalagHandlersTx(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		catalagHandlers := di.Get(ctx, constants.CatalogHandlersTxKey).(ddd.EventHandler[ddd.Event])

		return catalagHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])
	RegisterCatalagHandlers(subscriber, handlers)
}

func (h catalagHandlers[T]) HandleEvent(ctx context.Context, event T) (err error){
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling mall event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled catalog event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling catalog event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName(){
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	}
	return nil
}

func(h catalagHandlers[T])onProductAdded(ctx context.Context, event ddd.Event)error{
	payload := event.Payload().(*domain.Product)
	return h.catalog.AddProduct(ctx,payload.ID(),payload.StoreID,payload.Name,payload.Description,payload.SKU,payload.Price)
}
