package customerspb

import (
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
)

const (
	CustomerAggregateChannel = "mallbots.customers.events.Customer"

	CustomerRegisteredEvent = "customerspb.CustomerRegisteredEvent"

	CommandChannel = "mallbots.customers.Commands"
)

func (*CustomerRegistered) Key() string { return CustomerRegisteredEvent }

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) (err error) {
	registers := []registry.Registrable{
		// Customer events
		&CustomerRegistered{},
	}

	for _, item := range registers {
		err = serde.Register(item)
		if err != nil {
			return err
		}
	}

	return nil
}
