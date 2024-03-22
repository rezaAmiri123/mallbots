package am

import (
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	MessageBase interface {
		ddd.IDer
		Subject() string
		MessageName() string
		Metadata() ddd.Metadata
		SentAt() time.Time
	}
	Message interface {
		MessageBase
		Data() []byte
	}
	
	IncomingMessageBase interface {
		MessageBase
		ReceivedAt() time.Time
		Ack() error
		NAck() error
		Extend() error
		Kill() error
	}
	IncomingMessage interface {
		IncomingMessageBase
		Data() []byte
	}

	message struct {
		id       string
		name     string
		subject  string
		data     []byte
		metadata ddd.Metadata
		sentAt   time.Time
	}
)

var _ Message = (*message)(nil)

func (m message) ID() string             { return m.id }
func (m message) Subject() string        { return m.subject }
func (m message) MessageName() string    { return m.name }
func (m message) Metadata() ddd.Metadata { return m.metadata }
func (m message) SentAt() time.Time      { return m.sentAt }
func (m message) Data() []byte           { return m.data }
