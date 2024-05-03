package domain

import (
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/rezaAmiri123/edatV2/registry"
)

func RegistrationsWithSerde(serde registry.Serde) error {
	// Order
	if err := serde.Register(Order{}, func(v interface{}) error {
		basket := v.(*Order)
		basket.Aggregate = es.NewAggregate("", OrderAggregate)
		return nil
	}); err != nil {
		return err
	}

	// basket events
	if err := serde.Register(OrderCreated{}); err != nil {
		return err
	}


	// order snapshots
	if err := serde.RegisterKey(OrderV1{}.SnapshotName(), OrderV1{}); err != nil {
		return err
	}

	return nil
}
