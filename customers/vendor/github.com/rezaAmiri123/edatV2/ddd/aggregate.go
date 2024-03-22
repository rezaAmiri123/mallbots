package ddd

const (
	AggregateNameKey    = "aggregate-name"
	AggregateIDKey      = "aggregate-id"
	AggregateVersionKey = "aggregate-version"
)

type (
	AggregateIDer interface {
		AggregateID() string
	}
	AggregateNamer interface {
		AggregateName() string
	}
	AggregateVersioner interface {
		AggregateVersion() int
	}

	AggregateEvent interface {
		Event
		AggregateIDer
		AggregateNamer
		AggregateVersioner
	}
	Eventer interface {
		AddEvent(string, EventPayload, ...EventOption)
		Events() []AggregateEvent
		ClearEvents()
	}
	Aggregate interface {
		IDer
		AggregateNamer
		Eventer
		IDSetter
		NameSetter
	}

	aggregateEvent struct {
		event
	}
	aggregate struct {
		Entity
		events []AggregateEvent
	}
)

var _ interface {
	AggregateIDer
	AggregateNamer
	AggregateVersioner
} = (*aggregateEvent)(nil)

var _ Aggregate = (*aggregate)(nil)

func NewAggregate(id, name string) *aggregate {
	return &aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (e aggregateEvent) AggregateID() string   { return e.metadata.Get(AggregateIDKey).(string) }
func (e aggregateEvent) AggregateName() string { return e.metadata.Get(AggregateNameKey).(string) }
func (e aggregateEvent) AggregateVersion() int { return e.metadata.Get(AggregateVersionKey).(int) }

func (a aggregate) AggregateName() string     { return a.EntityName() }
func (a *aggregate) ClearEvents()             { a.events = []AggregateEvent{} }
func (a *aggregate) Events() []AggregateEvent { return a.events }
func (a aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(options, Metadata{
		AggregateIDKey:   a.ID(),
		AggregateNameKey: a.EntityName(),
	})
	a.events = append(a.events, aggregateEvent{
		event: newEvent(name, payload, options...),
	})
}

func (a *aggregate) setEvents(events []AggregateEvent) { a.events = events }
