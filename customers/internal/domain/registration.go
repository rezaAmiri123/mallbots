package domain

import "github.com/rezaAmiri123/edatV2/registry"

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
