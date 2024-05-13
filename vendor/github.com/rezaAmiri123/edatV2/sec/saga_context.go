package sec

type SagaContext[T any] struct {
	ID           string
	Data         T
	Step         int
	Done         bool
	Compensating bool
}

func (s *SagaContext[T]) advance(steps int) {
	var dir = 1
	if s.Compensating {
		dir = -1
	}

	s.Step += dir * steps
}

func (s *SagaContext[T]) complete() {
	s.Done = true
}

func (s *SagaContext[T]) compensate() {
	s.Compensating = true
}
