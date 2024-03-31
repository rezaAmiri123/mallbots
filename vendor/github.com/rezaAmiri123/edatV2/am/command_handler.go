package am

import (
	"context"
	"strings"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/registry"
)

type commandMsgHandler struct {
	reg        registry.Registry
	serializer MessageSerializer
	publisher  ReplyPublisher
	handler    ddd.CommandHandler[ddd.Command]
}

var _ MessageHandler = (*commandMsgHandler)(nil)

func NewCommandHandler(
	reg registry.Registry,
	serializer MessageSerializer,
	publisher ReplyPublisher,
	handler ddd.CommandHandler[ddd.Command],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(commandMsgHandler{
		reg:        reg,
		serializer: serializer,
		publisher:  publisher,
		handler:    handler,
	}, mws...)
}

func (h commandMsgHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	commandData, err := h.serializer.Deserialize(msg.Data())
	if err != nil {
		return err
	}

	commandName := msg.MessageName()
	payload, err := h.reg.Deserialize(commandName, commandData.Payload)
	if err != nil {
		return err
	}

	commandMsg := commandMessage{
		id:         msg.ID(),
		name:       commandName,
		payload:    payload,
		occurredAt: commandData.OccurredAt,
		msg:        msg,
	}

	destination := commandMsg.Metadata().Get(CommandReplyChannelHdr).(string)

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, destination, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, destination, h.success(reply, commandMsg))

}

func (h commandMsgHandler) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	return h.publisher.Publish(ctx, destination, reply)
}

func (h commandMsgHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeFailure)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHdr, OutcomeSuccess)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMsgHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHdr {
			continue
		}

		if strings.HasPrefix(key, CommandHdrPrefix) {
			hdr := ReplyHdrPrefix + key[len(CommandHdrPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}
	return reply
}
