package am

import (
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	ReplyMessage interface {
		MessageBase
		ddd.Reply
	}
	IncomingReplyMessage interface {
		IncomingMessageBase
		ddd.Reply
	}
	// ReplyPublisher  = MessagePublisher[ddd.Reply]
	// ReplySubscriber = MessageSubscriber[IncomingReplyMessage]
	// ReplyStream     = MessageStream[ddd.Reply, IncomingReplyMessage]

	replyMessage struct {
		id         string
		name       string
		payload    ddd.ReplyPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var _ IncomingReplyMessage = (*replyMessage)(nil)

func (r replyMessage) ID() string                { return r.id }
func (r replyMessage) ReplyName() string         { return r.name }
func (r replyMessage) Payload() ddd.ReplyPayload { return r.payload }
func (r replyMessage) Metadata() ddd.Metadata    { return r.msg.Metadata() }
func (r replyMessage) OccurredAt() time.Time     { return r.occurredAt }
func (r replyMessage) Subject() string           { return r.msg.Subject() }
func (r replyMessage) MessageName() string       { return r.msg.MessageName() }
func (r replyMessage) SentAt() time.Time         { return r.msg.SentAt() }
func (r replyMessage) ReceivedAt() time.Time     { return r.msg.ReceivedAt() }
func (r replyMessage) Ack() error                { return r.msg.Ack() }
func (r replyMessage) NAck() error               { return r.msg.NAck() }
func (r replyMessage) Extend() error             { return r.msg.Extend() }
func (r replyMessage) Kill() error               { return r.msg.Kill() }
