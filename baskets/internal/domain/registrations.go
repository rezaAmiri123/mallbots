package domain

import (
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/rezaAmiri123/edatV2/registry"
)

func RegistrationsWithSerde(serde registry.Serde) error {
	// Basket
	if err := serde.Register(Basket{}, func(v interface{}) error {
		basket := v.(*Basket)
		basket.Aggregate = es.NewAggregate("", BasketAggregate)
		basket.Items = make(map[string]Item)
		return nil
	}); err != nil {
		return err
	}

	// basket events
	if err := serde.Register(BasketStarted{}); err != nil {
		return err
	}

	return nil
}
