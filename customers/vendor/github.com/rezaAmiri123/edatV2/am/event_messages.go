package am

import (
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	EventMessage interface {
		MessageBase
		ddd.Event
	}
	IncomingEventMessage interface {
		IncomingMessageBase
		ddd.Event
	}

	eventMessage struct {
		id         string
		name       string
		payload    ddd.EventPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var _ IncomingEventMessage = (*eventMessage)(nil)

func (e eventMessage) ID() string                { return e.id}
func (e eventMessage) EventName() string         { return e.name}
func (e eventMessage) Payload() ddd.EventPayload { return e.payload}
func (e eventMessage) Metadata() ddd.Metadata    { return e.msg.Metadata()}
func (e eventMessage) OccurredAt() time.Time     { return e.occurredAt}
func (e eventMessage) Subject() string           { return e.msg.Subject()}
func (e eventMessage) MessageName() string       { return e.msg.MessageName()}
func (e eventMessage) SentAt() time.Time         { return e.msg.SentAt()}
func (e eventMessage) ReceivedAt() time.Time     { return e.msg.ReceivedAt()}
func (e eventMessage) Ack() error                { return e.msg.Ack()}
func (e eventMessage) NAck() error               { return e.msg.NAck()}
func (e eventMessage) Extend() error             { return e.msg.Extend()}
func (e eventMessage) Kill() error               { return e.msg.Kill()}
