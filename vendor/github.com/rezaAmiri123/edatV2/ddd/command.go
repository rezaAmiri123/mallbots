package ddd

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type (
	CommandHandler[T Command]interface{
		HandleCommand(ctx context.Context, cmd T)(Reply, error)
	}

	CommandHandlerFunc[T Command]func (ctx context.Context, cmd T)(Reply, error)


	CommandOption interface {
		configureCommand(*command)
	}

	CommandPayload any

	Command interface {
		IDer
		CommandName() string
		Payload() CommandPayload
		Metadata() Metadata
		OccurredAt() time.Time
	}

	command struct {
		Entity
		payload    CommandPayload
		metadata   Metadata
		occurredAt time.Time
	}
)

var _ Command = (*command)(nil)

func NewCommand(name string, payload CommandPayload, options ...CommandOption) Command {
	return newCommand(name, payload, options...)
}

func newCommand(name string, payload CommandPayload, options ...CommandOption) command {
	cmd := command{
		Entity:     NewEntity(uuid.New().String(), name),
		payload:    payload,
		metadata:   make(Metadata),
		occurredAt: time.Now(),
	}

	for _, option := range options {
		option.configureCommand(&cmd)
	}

	return cmd
}

func (c command) CommandName() string     { return c.EntityName() }
func (c command) Payload() CommandPayload { return c.payload }
func (c command) Metadata() Metadata      { return c.metadata }
func (c command) OccurredAt() time.Time   { return c.occurredAt }


func(f CommandHandlerFunc[T])HandleCommand(ctx context.Context, cmd T)(Reply, error){
	return f(ctx,cmd)
}