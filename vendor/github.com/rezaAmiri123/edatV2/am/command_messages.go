package am

import (
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	CommandMessage interface {
		MessageBase
		ddd.Command
	}
	IncomingCommandMessage interface {
		IncomingMessageBase
		ddd.Command
	}

	commandMessage struct {
		id         string
		name       string
		payload    ddd.CommandPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}
)

var _ IncomingCommandMessage = (*commandMessage)(nil)

func (c commandMessage) ID() string                  { return c.id }
func (c commandMessage) CommandName() string         { return c.name }
func (c commandMessage) Payload() ddd.CommandPayload { return c.payload }
func (c commandMessage) OccurredAt() time.Time       { return c.occurredAt }

func (c commandMessage) Subject() string        { return c.msg.Subject() }
func (c commandMessage) MessageName() string    { return c.msg.MessageName() }
func (c commandMessage) Metadata() ddd.Metadata { return c.msg.Metadata() }
func (c commandMessage) SentAt() time.Time      { return c.msg.SentAt() }
func (c commandMessage) ReceivedAt() time.Time  { return c.msg.ReceivedAt() }
func (c commandMessage) Ack() error             { return c.msg.Ack() }
func (c commandMessage) NAck() error            { return c.msg.NAck() }
func (c commandMessage) Extend() error          { return c.msg.Extend() }
func (c commandMessage) Kill() error            { return c.msg.Kill() }
