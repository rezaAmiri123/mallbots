package agent

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/mallbots/customers/internal/adapters"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
)

func (a *Agent) setupApplication() error {
	a.container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})

	a.container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresCustomerRepository(
			constants.CustomersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})

	// Prometheus counters
	customersRegistered := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.CustomersRegisteredCount,
	})

	// setup application
	a.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.NewInstrumentedApp(application.NewApplication(
			c.Get(constants.CustomersRepoKey).(domain.CustomerRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent]),
		), customersRegistered), nil
	})

	return nil
}
