package sec

type StepOption[T any] func(step *sagaStep[T])

func WithAction[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[notCompensating] = fn
	}
}

func WithCompensation[T any](fn StepActionFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.actions[isCompensating] = fn
	}
}

func OnActionReply[T any](replayName string, fn StepReplyHandlerFunc[T]) StepOption[T] {
	return func(step *sagaStep[T]) {
		step.handlers[notCompensating][replayName] = fn
	}
}

func InCompensationReply[T any](replyName string, fn StepReplyHandlerFunc[T])StepOption[T]{
	return func(step *sagaStep[T]) {
		step.handlers[isCompensating][replyName] = fn
	}
}
