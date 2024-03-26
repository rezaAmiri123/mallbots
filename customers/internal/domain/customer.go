package domain

import (
	"fmt"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/stackus/errors"
)

const CustomerAggregate = "customers.CustomerAggregate"

var (
	ErrNameCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the customer name cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the sms number cannot be blank")
)

type Customer struct {
	ddd.Aggregate
	Name      string
	SmsNumber string
	Enabled   bool
}

func (Customer) Key() string { return CustomerAggregate }

func NewCustomer(id string) *Customer {
	return &Customer{
		Aggregate: ddd.NewAggregate(id, CustomerAggregate),
	}
}

func RegisterCustomer(id, name, smsNumber string) (*Customer, error) {
	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}
	if name == "" {
		return nil, ErrNameCannotBeBlank
	}
	if smsNumber == "" {
		return nil, ErrSmsNumberCannotBeBlank
	}

	customer := NewCustomer(id)
	customer.Name = name
	customer.SmsNumber = smsNumber
	customer.Enabled = true

	customer.AddEvent(CustomerRegisteredEvent, &CustomerRegistered{
		Customer: customer,
	})
	fmt.Println(customer.Aggregate.Events())
	fmt.Println("customer.Events(): ",customer.Events())

	return customer, nil
}
