package events

import (
	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/sec"
	"github.com/rezaAmiri123/mallbots/cosec/internal"
	"github.com/rezaAmiri123/mallbots/cosec/internal/models"
)

func NewReplyHandlers(
	reg registry.Registry,
	serializer am.MessageSerializer,
	orchestrator sec.Orchestrator[*models.CreateOrderData],
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewReplyHandler(reg, serializer, orchestrator, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(
		internal.CreateOrderReplyChannel,
		handlers,
		am.GroupName("cosec-replies"),
	)
	return err
}
