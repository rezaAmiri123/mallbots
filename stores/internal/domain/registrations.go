package domain

import (
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/rezaAmiri123/edatV2/registry"
)

func Registrations(serde registry.Serde) (err error) {
	if err = serde.Register(Store{}, func(v interface{}) error {
		store := v.(*Store)
		store.Aggregate = es.NewAggregate("", StoreAggregate)
		return nil
	}); err != nil {
		return
	}

	// store events
	if err = serde.Register(StoreCreated{}); err != nil {
		return
	}

	// store snapshots
	if err = serde.RegisterKey(StoreV1{}.SnapshotName(), StoreV1{}); err != nil {
		return
	}

	if err := serde.Register(Product{}, func(v interface{}) error {
		product := v.(*Product)
		product.Aggregate = es.NewAggregate("", ProductAggregate)
		return nil
	}); err != nil {
		return err
	}

	// product events
	if err = serde.Register(ProductAdded{}); err != nil {
		return
	}
	if err = serde.RegisterKey(ProductPriceIncreasedEvent, ProductPriceChanged{}); err != nil {
		return
	}

	// product snapshots
	if err = serde.RegisterKey(ProductV1{}.SnapshotName(), ProductV1{}); err != nil {
		return
	}

	return
}
