package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/ordering/internal/domain"
	"github.com/rezaAmiri123/mallbots/ordering/orderingpb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.OrderCreatedEvent,
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
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*domain.Order)
	items := make([]*orderingpb.OrderCreated_Item, len(payload.Items))
	for i, item := range payload.Items {
		items[i] = &orderingpb.OrderCreated_Item{
			StoreId:   item.StoreID,
			ProductId: item.ProductID,
			Price:     item.Price,
			Quantity:  int32(item.Quantity),
		}
	}

	newEvent := ddd.NewEvent(orderingpb.OrderCreatedEvent, &orderingpb.OrderCreated{
		Id:         payload.ID(),
		CustomerId: payload.CustomerID,
		PaymentId:  payload.PaymentID,
		ShoppingId: payload.ShoppingID,
		Items:      items,
	})

	return h.publisher.Publish(ctx, orderingpb.OrderAggregateChannel, newEvent)
}
