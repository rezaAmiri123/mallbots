package domain

import (
	"github.com/rezaAmiri123/edatV2/registry"
)

func RegistrationsWithSerde(serde registry.Serde) error {
	// Events
	// it is not event source, so we do not need to register internal events
	return nil
}
