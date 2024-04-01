package am

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type eventMsgHandler struct {
	reg        registry.Registry
	serializer MessageSerializer
	handler    ddd.EventHandler[ddd.Event]
}

func NewEventHandler(
	reg registry.Registry,
	serializer MessageSerializer,
	handler ddd.EventHandler[ddd.Event],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(eventMsgHandler{
		reg:        reg,
		serializer: serializer,
		handler:    handler,
	}, mws...)
}

var _ MessageHandler = (*eventMsgHandler)(nil)

func(h eventMsgHandler)HandleMessage(ctx context.Context, msg IncomingMessage) error{
	eventData, err := h.serializer.Deserialize(msg.Data())
	if err != nil{
		return err
	}

	eventName := msg.MessageName()

	payload, err := h.reg.Deserialize(eventName, eventData.Payload)
	if err != nil{
		return err
	}

	// TODO either this should be a ddd.Event or the handler is a HandleMessage[am.EventMessage]
	eventMsg := eventMessage{
		id: msg.ID(),
		name: eventName,
		payload: payload,
		occurredAt: eventData.OccurredAt,
		msg: msg,
	}

	return h.handler.HandleEvent(ctx,eventMsg)
}
