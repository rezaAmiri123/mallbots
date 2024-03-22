// Code generated by mockery v2.33.0. DO NOT EDIT.

package ddd

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockEventHandler is an autogenerated mock type for the EventHandler type
type MockEventHandler[T Event] struct {
	mock.Mock
}

// HandleEvent provides a mock function with given fields: ctx, event
func (_m *MockEventHandler[T]) HandleEvent(ctx context.Context, event T) error {
	ret := _m.Called(ctx, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, T) error); ok {
		r0 = rf(ctx, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockEventHandler creates a new instance of MockEventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventHandler[T Event](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventHandler[T] {
	mock := &MockEventHandler[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
