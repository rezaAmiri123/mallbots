package am

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type (
	replyHandler struct {
		reg        registry.Registry
		serializer MessageSerializer
		handler    ddd.ReplyHandler[ddd.Reply]
	}
)

var _ MessageHandler = (*replyHandler)(nil)

func NewReplyHandler(
	reg registry.Registry,
	serializer MessageSerializer,
	handler ddd.ReplyHandler[ddd.Reply],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(replyHandler{
		reg:        reg,
		serializer: serializer,
		handler:    handler,
	}, mws...)
}

func (h replyHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	replyData, err := h.serializer.Deserialize(msg.Data())

	if err != nil {
		return err
	}

	replyName := msg.MessageName()
	var payload any

	if replyName != SuccessReply && replyName != FailureReply {
		payload, err = h.reg.Deserialize(replyName, replyData.Payload)
		if err != nil {
			return nil
		}
	}

	replyMsg := replyMessage{
		id:         msg.ID(),
		name:       replyName,
		payload:    payload,
		occurredAt: replyData.OccurredAt,
		msg:        msg,
	}

	return h.handler.HandleReply(ctx, replyMsg)
}
