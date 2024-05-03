package basketspb

import "github.com/rezaAmiri123/edatV2/registry"

const (
	BasketAggregateChannel = "mallbots.baskets.events.Basket"

	BasketStartedEvent    = "basketsapi.BasketStarted"
	BasketCheckedOutEvent = "basketsapi.BasketCheckedOut"
)

func RegistrationsWithSerde(serde registry.Serde) error {
	// Basket events
	if err := serde.Register(&BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCheckedOut{}); err != nil {
		return err
	}

	return nil
}

func (*BasketStarted) Key() string    { return BasketStartedEvent }
func (*BasketCheckedOut) Key() string { return BasketCheckedOutEvent }
