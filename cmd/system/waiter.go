package system

import "github.com/rezaAmiri123/edatV2/waiter"

type Waiter interface{
	Waiter() waiter.Waiter
}

func (s *System) initWaiter() {
	s.waiter = waiter.New(waiter.CatchSignals())
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}
