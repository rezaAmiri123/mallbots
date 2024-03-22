package am

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/stackus/errors"
)

type (
	fakeEventMessage struct {
		subject string
		payload ddd.Event
	}

	FakeEventPublisher struct {
		messages []fakeEventMessage
	}
)

var _ EventPublisher = (*FakeEventPublisher)(nil)

func NewFakeEventPublisher() *FakeEventPublisher {
	return &FakeEventPublisher{
		messages: []fakeEventMessage{},
	}
}

func (p *FakeEventPublisher) Publish(ctx context.Context, topicName string, event ddd.Event) error {
	p.messages = append(p.messages, fakeEventMessage{
		subject: topicName,
		payload: event,
	})

	return nil
}

func (p *FakeEventPublisher) Reset() {
	p.messages = []fakeEventMessage{}
}

func (p *FakeEventPublisher) Last() (string, ddd.Event, error) {
	var v ddd.Event
	if len(p.messages) == 0 {
		return "", v, errors.ErrNotFound.Msg("no events have been published")
	}

	last := p.messages[len(p.messages)-1]
	
	return last.subject, last.payload, nil
}
