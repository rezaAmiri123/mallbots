package agent

import (
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
)

func (a *Agent) setupRegistry() error {
	reg := registry.NewRegistry()
	serde := serdes.NewJsonSerde(reg)

	domain.RegistrationsWithSerde(serde)

	a.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		return reg, nil
	})
	return nil
}
