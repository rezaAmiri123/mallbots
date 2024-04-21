package events

import (
	"context"
	"time"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/errorsotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type integrationHandlers[T ddd.Event] struct {
	stores domain.StoreCacheRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationHandlers(
	reg registry.Registry,
	serializer am.MessageSerializer,
	stores domain.StoreCacheRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg,serializer,integrationHandlers[ddd.Event]{
		stores: stores,
	},mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler)(err error){
	_, err = subscriber.Subscribe(storespb.StoreAggregateChannel,handlers,am.MessageFilter{
		storespb.StoreCreatedEvent,
	}, am.GroupName("baskets-stores"))
	
	return err
}

func(h integrationHandlers[T])HandleEvent(ctx context.Context, event T) (err error){
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

	switch event.EventName(){
	case storespb.StoreCreatedEvent:
		return h.onStoreCreated(ctx,event)
	}

	return nil
}
func(h integrationHandlers[T])onStoreCreated(ctx context.Context, event ddd.Event)error{
	payload := event.Payload().(*storespb.StoreCreated)
	return h.stores.Add(ctx, payload.Id,payload.Name)
}
// func(h integrationHandlers[T]){}
// func(h integrationHandlers[T]){}
// func(h integrationHandlers[T]){}