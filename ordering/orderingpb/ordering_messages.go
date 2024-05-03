package orderingpb

import "github.com/rezaAmiri123/edatV2/registry"

const (
	OrderAggregateChannel = "mallbots.ordering.events.Order"

	OrderCreatedEvent = "ordersapi.OrderCreated"
)

func RegistrationsWithSerde(serde registry.Serde) error {
	// order events
	if err := serde.Register(&OrderCreated{}); err != nil {
		return err
	}

	return nil
}

func (*OrderCreated) Key() string { return OrderCreatedEvent }
