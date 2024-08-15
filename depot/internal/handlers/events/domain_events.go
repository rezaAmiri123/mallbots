package events

// import (
// 	"github.com/rezaAmiri123/edatV2/am"
// 	"github.com/rezaAmiri123/edatV2/ddd"
// 	"github.com/rezaAmiri123/mallbots/depot/internal/domain"
// )

// type domainHandlers[T ddd.AggregateEvent] struct{
// 	publisher am.EventPublisher
// }

// var _ ddd.EventHandler[ddd.AggregateEvent]=(*domainHandlers[ddd.AggregateEvent])(nil)

// func NewDomainEventHandlers(publisher am.EventPublisher)ddd.EventHandler[ddd.AggregateEvent]{
// 	return domainHandlers[ddd.AggregateEvent]{
// 		publisher: publisher,
// 	}
// }

// func RegisterDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.AggregateEvent], handlers ddd.EventHandler[ddd.AggregateEvent]){
// 	subscriber.Subscribe(handlers,domain.ShoppingListC)
// }
