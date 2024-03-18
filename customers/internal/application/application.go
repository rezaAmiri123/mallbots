package application

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
)

type (
	RegisterCustomer struct {
		ID        string
		Name      string
		SmsNumber string
	}
)

type (
	App interface {
		RegisterCustomer(ctx context.Context, register RegisterCustomer) error
	}

	Application struct {
		customers       domain.CustomerRepository
		domainPublisher ddd.EventPublisher[ddd.AggregateEvent]
	}
)

var _ App = (*Application)(nil)

func NewApplication(
	customers domain.CustomerRepository,
	domainPublisher ddd.EventPublisher[ddd.AggregateEvent],
) *Application {
	return &Application{
		customers:       customers,
		domainPublisher: domainPublisher,
	}
}

func (a Application) RegisterCustomer(ctx context.Context, register RegisterCustomer) error {
	customer, err := domain.RegisterCustomer(register.ID, register.Name, register.SmsNumber)
	if err != nil {
		return err
	}

	if err = a.customers.Save(ctx, customer); err != nil {
		return err
	}

	// publish domain events
	if err = a.domainPublisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}
