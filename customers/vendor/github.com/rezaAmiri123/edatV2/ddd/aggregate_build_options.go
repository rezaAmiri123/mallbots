package ddd

import (
	"fmt"

	"github.com/rezaAmiri123/edatV2/registry"
)

type EventsSetter interface {
	setEvents([]Event)
}

func SetEvents(events ...Event) registry.BuildOption {
	return func(v interface{}) error {
		if agg, ok := v.(EventsSetter); ok {
			agg.setEvents(events)
			return nil
		}
		return fmt.Errorf("%T does not have the method setEvents([]ddd.Event)")
	}
}
