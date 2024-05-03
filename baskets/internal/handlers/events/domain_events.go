package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/mallbots/baskets/basketspb"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
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
		domain.BasketStartedEvent,
		domain.BasketCheckedOutEvent,
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
	case domain.BasketStartedEvent:
		return h.onBasketStarted(ctx, event)
	case domain.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)

	}

	return nil
}

func (h domainHandlers[T]) onBasketStarted(ctx context.Context, event ddd.Event) (err error) {
	basket := event.Payload().(*domain.Basket)

	newEvent := ddd.NewEvent(basketspb.BasketStartedEvent, &basketspb.BasketStarted{
		Id:         basket.ID(),
		CustomerId: basket.CustomrID,
	})

	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel, newEvent)
}

func (h domainHandlers[T]) onBasketCheckedOut(ctx context.Context, event ddd.Event) (err error) {
	basket := event.Payload().(*domain.Basket)
	items := make([]*basketspb.BasketCheckedOut_Item, 0, len(basket.Items))
	for _, item := range basket.Items {
		items = append(items, &basketspb.BasketCheckedOut_Item{
			StoreId:     item.StoreID,
			StoreName:   item.StoreName,
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.ProductPrice,
			Quantity:    int32(item.Quantity),
		})
	}

	newEvent := ddd.NewEvent(basketspb.BasketCheckedOutEvent, &basketspb.BasketCheckedOut{
		Id:         basket.ID(),
		CustomerId: basket.CustomrID,
		PaymentId:  basket.PaymentID,
		Items:      items,
	})

	return h.publisher.Publish(ctx, basketspb.BasketAggregateChannel, newEvent)
}
