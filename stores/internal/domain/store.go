package domain

import (
	"fmt"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/stackus/errors"
)

const StoreAggregate = "stores.Store"

var (
	ErrStoreNameIsBlank               = errors.Wrap(errors.ErrBadRequest, "the store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.Wrap(errors.ErrBadRequest, "the store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.Wrap(errors.ErrBadRequest, "the store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.Wrap(errors.ErrBadRequest, "the store is already not participating")
)

type Store struct {
	es.Aggregate
	Name          string
	Location      string
	Participating bool
}

func (Store) Key() string { return StoreAggregate }

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Store)(nil)

func NewStore(id string) *Store {
	return &Store{
		Aggregate: es.NewAggregate(id, StoreAggregate),
	}
}

func (s *Store) InitStore(name, location string) (ddd.Event, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	s.AddEvent(StoreCreatedEvent, &StoreCreated{
		Name:     name,
		Location: location,
	})

	return ddd.NewEvent(StoreCreatedEvent, s), nil
}
func (s *Store) ApplyEvent(event ddd.Event) error {
	fmt.Println("func (s *Store) ApplyEvent(event ddd.Event) error {")
	switch payload := event.Payload().(type) {
	case *StoreCreated:
		s.Name = payload.Name
		s.Location = payload.Location
	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}
// ApplySnapshot implements es.Snapshotter
func (s *Store) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *StoreV1:
		s.Name = ss.Name
		s.Location = ss.Location
		s.Participating = ss.Participating

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", s, snapshot)
	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s Store) ToSnapshot() es.Snapshot {
	return StoreV1{
		Name:          s.Name,
		Location:      s.Location,
		Participating: s.Participating,
	}
}
