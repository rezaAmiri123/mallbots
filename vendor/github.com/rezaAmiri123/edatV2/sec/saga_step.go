package sec

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	SagaStep[T any] interface {
		Action(fn StepActionFunc[T]) SagaStep[T]
		Compensation(fn StepActionFunc[T]) SagaStep[T]
		OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]
		OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T]

		isInvocable(compensating bool) bool
		execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T]
		handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error
	}

	sagaStep[T any] struct {
		actions  map[bool]StepActionFunc[T]
		handlers map[bool]map[string]StepReplyHandlerFunc[T]
	}
)

var _ SagaStep[any] = (*sagaStep[any])(nil)

func (s *sagaStep[T]) Action(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[notCompensating] = fn
	return s
}

func (s *sagaStep[T]) Compensation(fn StepActionFunc[T]) SagaStep[T] {
	s.actions[isCompensating] = fn
	return s
}

func (s *sagaStep[T]) OnActionReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[notCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) OnCompensationReply(replyName string, fn StepReplyHandlerFunc[T]) SagaStep[T] {
	s.handlers[isCompensating][replyName] = fn
	return s
}

func (s *sagaStep[T]) isInvocable(compensating bool) bool {
	return s.actions[compensating] != nil
}

func (s *sagaStep[T]) execute(ctx context.Context, sagaCtx *SagaContext[T]) stepResult[T] {
	if action := s.actions[sagaCtx.Compensating]; action != nil {
		destination, cmd, err := action(ctx, sagaCtx.Data)
		return stepResult[T]{
			ctx:         sagaCtx,
			destination: destination,
			cmd:         cmd,
			err:         err,
		}
	}
	return stepResult[T]{ctx: sagaCtx}
}

func (s *sagaStep[T]) handle(ctx context.Context, sagaCtx *SagaContext[T], reply ddd.Reply) error {
	if handler := s.handlers[sagaCtx.Compensating][reply.ReplyName()]; handler != nil {
		return handler(ctx, sagaCtx.Data, reply)
	}
	return nil
}
