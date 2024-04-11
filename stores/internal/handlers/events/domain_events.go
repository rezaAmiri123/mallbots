package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.Event] = (*domainHandlers[ddd.Event])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.StoreCreatedEvent,
		domain.ProductAddedEvent,
		domain.ProductPriceIncreasedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling domain event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled domain event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling domain event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	case domain.ProductPriceIncreasedEvent:
		return h.onProductPriceIncreased(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onStoreCreated(ctx context.Context, event ddd.Event) error {
	store := event.Payload().(*domain.Store)

	newEvent := ddd.NewEvent(storespb.StoreCreatedEvent, &storespb.StoreCreated{
		Id:       store.ID(),
		Name:     store.Name,
		Location: store.Location,
	})

	return h.publisher.Publish(ctx, storespb.StoreAggregateChannel, newEvent)
}

func (h domainHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.Product)

	newEvent := ddd.NewEvent(storespb.ProductAddedEvent, &storespb.ProductAdded{
		Id:          product.ID(),
		StoreId:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.SKU,
		Price:       product.Price,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, newEvent)
}

func (h domainHandlers[T]) onProductPriceIncreased(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.ProductPriceDelta)

	newEvent := ddd.NewEvent(storespb.ProductPriceIncreasedEvent, &storespb.ProductPriceChanged{
		Id:    payload.Product.ID(),
		Delta: payload.Delta,
	})

	return h.publisher.Publish(ctx, storespb.ProductAggregateChannel, newEvent)
}
