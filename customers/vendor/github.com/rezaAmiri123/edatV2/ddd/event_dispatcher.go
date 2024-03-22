package ddd

import (
	"context"
	"sync"
)

type (
	EventSubscriber[T Event] interface {
		Subscribe(handler EventHandler[T], events ...string)
	}
	EventPublisher[T Event] interface {
		Publish(ctx context.Context, events ...T) error
	}

	eventHandler[T Event] struct {
		h       EventHandler[T]
		filters map[string]struct{}
	}

	EventDispatcher[T Event] struct {
		handlers []eventHandler[T]
		mu       sync.Mutex
	}
)

var _ interface {
	EventSubscriber[Event]
	EventPublisher[Event]
} = (*EventDispatcher[Event])(nil)

func NewEventDispatcher[T Event]() *EventDispatcher[T] {
	return &EventDispatcher[T]{
		handlers: make([]eventHandler[T], 0),
	}
}

func (h *EventDispatcher[T]) Publish(ctx context.Context, events ...T) error {
	for _, event := range events{
		for _, handler := range h.handlers{
			if handler.filters!=nil{
				if _,exists := handler.filters[event.EventName()];!exists{
					continue
				}
			}
			err := handler.h.HandleEvent(ctx,event)
			if err!= nil{
				return err
			}
		}
	}
	return nil
}

func (h *EventDispatcher[T]) Subscribe(handler EventHandler[T], events ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var filters map[string]struct{}
	if len(events) > 0 {
		filters = make(map[string]struct{})
		for _, event := range events {
			filters[event] = struct{}{}
		}
	}

	h.handlers = append(h.handlers, eventHandler[T]{
		h:       handler,
		filters: filters,
	})
}
