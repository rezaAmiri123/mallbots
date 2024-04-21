package domain

import (
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/stackus/errors"
)

var (
	ErrBasketHasNoItems         = errors.Wrap(errors.ErrBadRequest, "the basket has no items")
	ErrBasketCannotBeModified   = errors.Wrap(errors.ErrBadRequest, "the basket cannot be modified")
	ErrBasketCannotBeCancelled  = errors.Wrap(errors.ErrBadRequest, "the basket cannot be cancelled")
	ErrQuantityCannotBeNegative = errors.Wrap(errors.ErrBadRequest, "the item quantity cannot be negative")
	ErrBasketIDCannotBeBlank    = errors.Wrap(errors.ErrBadRequest, "the basket id cannot be blank")
	ErrPaymentIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrCustomerIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
)

const BasketAggregate = "baskets.Basket"

type Basket struct {
	es.Aggregate
	CustomrID string
	PaymentID string
	Items     map[string]Item
	Status    BasketStatus
}

func (Basket) Key() string { return BasketAggregate }

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Basket)(nil)

func NewBasket(id string) *Basket {
	return &Basket{
		Aggregate: es.NewAggregate(id, BasketAggregate),
		Items:     make(map[string]Item),
	}
}
func (b *Basket) Start(customerID string) (ddd.Event, error) {
	if b.Status != BasketUnknown {
		return nil, ErrBasketCannotBeModified
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	b.AddEvent(BasketStartedEvent, &BasketStarted{
		CustomerID: customerID,
	})

	return ddd.NewEvent(BasketStartedEvent, b), nil
}

// func (b *Basket){}
// func (b *Basket){}
// func (b *Basket){}
// func (b *Basket){}
func (b *Basket) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *BasketStarted:
		b.CustomrID = payload.CustomerID
	default:
		errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}
	return nil
}

func (b *Basket) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *BasketV1:
		b.CustomrID = ss.CustomerID
		b.PaymentID = ss.PaymentID
		b.Items = ss.Items
		b.Status = ss.Status
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *Basket) ToSnapshot() es.Snapshot {
	return &BasketV1{
		CustomerID: b.CustomrID,
		PaymentID:  b.PaymentID,
		Items:      b.Items,
		Status:     b.Status,
	}
}
