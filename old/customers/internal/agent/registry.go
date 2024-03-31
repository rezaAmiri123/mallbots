package agent

import (
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
)

func (a *Agent) setupRegistry() error {
	var serde registry.Serde
	reg := registry.NewRegistry()

	switch a.config.SerdeType {
	default:
		serde = serdes.NewJsonSerde(reg)
	}

	if err := domain.RegistrationsWithSerde(serde); err != nil {
		return err
	}

	if err := customerspb.RegistrationsWithSerde(serde); err != nil {
		return err
	}

	a.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		return reg, nil
	})
	return nil
}
