package sec

type (
	Saga[T any] interface {
		AddStep() SagaStep[T]
		Name() string
		ReplyTopic() string
		getSteps() []SagaStep[T]
	}

	saga[T any] struct {
		name       string
		replyTopic string
		steps      []SagaStep[T]
	}
)

var _ Saga[any] = (*saga[any])(nil)

func NewSaga[T any](name, replyTopic string) Saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

func (s *saga[T]) AddStep() SagaStep[T] {
	step := &sagaStep[T]{
		actions: map[bool]StepActionFunc[T]{
			notCompensating: nil,
			isCompensating:  nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFunc[T]{
			notCompensating: {},
			isCompensating:  {},
		},
	}
	s.steps = append(s.steps, step)

	return step
}

func (s *saga[T]) Name() string {
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga[T]) getSteps() []SagaStep[T] {
	return s.steps
}
