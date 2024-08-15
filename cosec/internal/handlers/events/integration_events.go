package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/sec"
	"github.com/rezaAmiri123/mallbots/cosec/internal/models"
	"github.com/rezaAmiri123/mallbots/ordering/orderingpb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	orchestrator sec.Orchestrator[*models.CreateOrderData]
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationHandlers(
	reg registry.Registry,
	serializer am.MessageSerializer,
	orchestrator sec.Orchestrator[*models.CreateOrderData],
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, serializer, integrationHandlers[ddd.Event]{
		orchestrator: orchestrator,
	}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) (err error) {
	_, err = subscriber.Subscribe(orderingpb.OrderAggregateChannel, handlers, am.MessageFilter{
		orderingpb.OrderCreatedEvent,
	}, am.GroupName("cosec-ordering"))

	return err
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) (err error) {
	span := trace.SpanFromContext(ctx)
	defer func(started time.Time) {
		if err != nil {
			span.AddEvent(
				"Encountered an error handling integration event",
				trace.WithAttributes(errorsotel.ErrAttrs(err)...),
			)
		}
		span.AddEvent("Handled integration event", trace.WithAttributes(
			attribute.Int64("TookMS", time.Since(started).Milliseconds()),
		))
	}(time.Now())

	span.AddEvent("Handling integration event", trace.WithAttributes(
		attribute.String("Event", event.EventName()),
	))

	switch event.EventName() {
	case orderingpb.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	}

	return nil
}
func (h integrationHandlers[T]) onOrderCreated(ctx context.Context, event ddd.Event) error {
	payload := event.Payload().(*orderingpb.OrderCreated)

	var total float64
	items := make([]models.Item, 0, len(payload.GetItems()))
	for _, item := range payload.GetItems() {
		items = append(items, models.Item{
			ProductID: item.ProductId,
			StoreID:   item.StoreId,
			Price:     item.Price,
			Quantity:  int(item.Quantity),
		})
		total += float64(item.Quantity) * item.Price
	}

	data := &models.CreateOrderData{
		OrderID: payload.GetId(),
		CustomerID: payload.GetCustomerId(),
		PaymentID: payload.GetPaymentId(),
		Items: items,
		Total: total,
	}

	// Start the CreateOrderSaga
	return h.orchestrator.Start(ctx,event.ID(),data)
}
