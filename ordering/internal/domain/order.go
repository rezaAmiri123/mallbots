package domain

import (
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/stackus/errors"
)

const OrderAggregate = "ordering.Order"

var (
	ErrOrderAlreadyCreated     = errors.Wrap(errors.ErrBadRequest, "the order cannot be recreated")
	ErrOrderHasNoItems         = errors.Wrap(errors.ErrBadRequest, "the order has no items")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "payment id cannot be blank")
)

type Order struct {
	es.Aggregate
	CustomerID string
	PaymentID  string
	InvoiceID  string
	ShoppingID string
	Items      []Item
	Status     OrderStatus
}

var _ interface {
	es.EventApplier
	es.Snapshotter
	registry.Registrable
} = (*Order)(nil)

func NewOrder(id string) *Order {
	return &Order{
		Aggregate: es.NewAggregate(id, OrderAggregate),
	}
}

func (o *Order) CreateOrder(id, customerID, paymentID string, items []Item) (ddd.Event, error) {
	if o.Status != OrderUnknown {
		return nil, ErrOrderAlreadyCreated
	}

	if len(items) == 0 {
		return nil, ErrOrderHasNoItems
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == "" {
		return nil, ErrPaymentIDCannotBeBlank
	}

	o.AddEvent(OrderCreatedEvent, &OrderCreated{
		CustomerID: o.CustomerID,
		PaymentID:  o.PaymentID,
		Items:      o.Items,
	})

	return ddd.NewEvent(OrderCreatedEvent, o), nil
}

// func(o *Order){}
// func(o *Order){}
// func(o *Order){}
// func(o *Order){}
// func(o *Order){}

func (Order) Key() string { return OrderAggregate }

func (o *Order) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *OrderCreated:
		o.CustomerID = payload.CustomerID
		o.PaymentID = payload.PaymentID
		o.Items = payload.Items
		o.Status = OrderIsPending
	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", o, event.EventName(), payload)
	}

	return nil
}

func (o *Order) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OrderV1:
		o.CustomerID = ss.CustomerID
		o.PaymentID = ss.PaymentID
		o.InvoiceID = ss.InvoiceID
		o.ShoppingID = ss.ShoppingID
		o.Items = ss.Items
		o.Status = ss.Status
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", o, snapshot)
	}

	return nil
}

func (o *Order) ToSnapshot() es.Snapshot {
	return &OrderV1{
		CustomerID: o.CustomerID,
		PaymentID:  o.PaymentID,
		InvoiceID:  o.InvoiceID,
		ShoppingID: o.ShoppingID,
		Items:      o.Items,
		Status:     o.Status,
	}
}
