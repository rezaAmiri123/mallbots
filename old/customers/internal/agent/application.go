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
	a.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
		customerRepo := adapters.NewPostgresCustomerRepository(
			constants.CustomersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		)
		return application.NewInstrumentedApp(application.NewApplication(
			customerRepo,
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent]),
		), customersRegistered), nil
	})
	// setup application
	a.container.AddSingleton(constants.ApplicationKey, func(c di.Container) (any, error) {
		customerRepo := adapters.NewPostgresCustomerRepository(
			constants.CustomersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
		)
		return application.NewInstrumentedApp(application.NewApplication(
			customerRepo,
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent]),
		), customersRegistered), nil
	})

	return nil
}

func (a *Agent) cleanupApplication() error {
	return nil
}